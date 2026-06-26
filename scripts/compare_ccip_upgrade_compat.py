#!/usr/bin/env python3
"""Compare deployed CCIP package JSON (normalized RPC) against local Move source."""

from __future__ import annotations

import argparse
import glob
import json
import re
import sys
import urllib.error
import urllib.request
from collections import Counter
from copy import deepcopy
from dataclasses import dataclass, field
from pathlib import Path
from typing import Any

REPO_ROOT = Path(__file__).resolve().parents[1]
SRC_GLOB = [
    "contracts/ccip/ccip/sources/*.move",
    "contracts/ccip/ccip/sources/util/*.move",
]

ADDR_MAP_LOCAL = {
    "std": "0x1",
    "sui": "0x2",
    "ccip": "0xSELF",
    "mcms": "0xMCMS",
    "fast_mcms": "0xMCMS",
}

IMPLICIT_NAME = {
    "TxContext": ("0x2", "tx_context", "TxContext"),
    "UID": ("0x2", "object", "UID"),
    "ID": ("0x2", "object", "ID"),
    "Option": ("0x1", "option", "Option"),
    "String": ("0x1", "string", "String"),
}

RPC_URLS = {
    "mainnet": [
        "https://fullnode.mainnet.sui.io:443",
        "https://mainnet.sui.rpcpool.com",
    ],
    "testnet": [
        "https://fullnode.testnet.sui.io:443",
    ],
}


@dataclass
class CompareReport:
    label: str
    issues: list[str] = field(default_factory=list)
    unparsed: list[str] = field(default_factory=list)
    additions_funs: list[str] = field(default_factory=list)
    additions_structs: list[str] = field(default_factory=list)
    stats: dict[str, Any] = field(default_factory=dict)


def fetch_normalized_modules(package_id: str, network: str) -> dict[str, Any]:
    payload = json.dumps(
        {
            "jsonrpc": "2.0",
            "id": 1,
            "method": "sui_getNormalizedMoveModulesByPackage",
            "params": [package_id],
        }
    ).encode()
    headers = {"Content-Type": "application/json"}
    errors: list[str] = []
    for url in RPC_URLS.get(network, RPC_URLS["mainnet"]):
        req = urllib.request.Request(url, data=payload, headers=headers, method="POST")
        opener = urllib.request.build_opener(urllib.request.ProxyHandler({}))
        try:
            with opener.open(req, timeout=120) as resp:
                data = json.loads(resp.read())
            if "error" in data:
                errors.append(f"{url}: {data['error']}")
                continue
            result = data.get("result")
            if not result:
                errors.append(f"{url}: empty result")
                continue
            return result
        except (urllib.error.URLError, TimeoutError, json.JSONDecodeError) as exc:
            errors.append(f"{url}: {exc}")
    raise RuntimeError("RPC fetch failed:\n  " + "\n  ".join(errors))


def fetch_upgrade_cap(package_id: str, upgrade_cap_id: str, network: str) -> dict[str, Any]:
    payload = json.dumps(
        {
            "jsonrpc": "2.0",
            "id": 1,
            "method": "sui_getObject",
            "params": [upgrade_cap_id, {"showContent": True}],
        }
    ).encode()
    headers = {"Content-Type": "application/json"}
    for url in RPC_URLS.get(network, RPC_URLS["mainnet"]):
        req = urllib.request.Request(url, data=payload, headers=headers, method="POST")
        opener = urllib.request.build_opener(urllib.request.ProxyHandler({}))
        try:
            with opener.open(req, timeout=60) as resp:
                data = json.loads(resp.read())
            fields = (
                data.get("result", {})
                .get("data", {})
                .get("content", {})
                .get("fields", {})
            )
            return {
                "upgrade_cap_id": upgrade_cap_id,
                "package_id": package_id,
                "version": fields.get("version"),
                "policy": fields.get("policy"),
            }
        except (urllib.error.URLError, TimeoutError, json.JSONDecodeError):
            continue
    return {"upgrade_cap_id": upgrade_cap_id, "package_id": package_id, "version": None, "policy": None}


def discover_addr_map(deployed: dict[str, Any]) -> dict[str, str]:
    counts: Counter[str] = Counter()

    def walk(x: Any) -> None:
        if isinstance(x, dict):
            if "Struct" in x:
                counts[x["Struct"]["address"].lower()] += 1
                for a in x["Struct"].get("typeArguments", []):
                    walk(a)
            elif "Reference" in x:
                walk(x["Reference"])
            elif "MutableReference" in x:
                walk(x["MutableReference"])
            elif "Vector" in x:
                walk(x["Vector"])
        elif isinstance(x, list):
            for it in x:
                walk(it)

    for mod in deployed.values():
        for s in mod.get("structs", {}).values():
            for f in s.get("fields", []):
                walk(f.get("type"))
        for fn in mod.get("exposedFunctions", {}).values():
            for p in fn.get("parameters", []):
                walk(p)
            for r in fn.get("return", []):
                walk(r)

    self_addr = counts.most_common(1)[0][0] if counts else ""
    external = [
        a
        for a in counts
        if a not in (self_addr, "0x1", "0x2")
    ]
    mcms = external[0] if external else ""
    return {
        self_addr: "0xSELF",
        mcms: "0xMCMS",
        "0x2": "0x2",
        "0x1": "0x1",
    }


def canon_deployed(t: Any, addr_map: dict[str, str], type_param_names: list[str]) -> str:
    if isinstance(t, str):
        return {
            "U8": "u8",
            "U16": "u16",
            "U32": "u32",
            "U64": "u64",
            "U128": "u128",
            "U256": "u256",
            "Bool": "bool",
            "Address": "address",
            "Signer": "signer",
        }.get(t, t)
    if isinstance(t, dict):
        if "TypeParameter" in t:
            idx = t["TypeParameter"]
            return type_param_names[idx] if idx < len(type_param_names) else f"T{idx}"
        if "Reference" in t:
            return "&" + canon_deployed(t["Reference"], addr_map, type_param_names)
        if "MutableReference" in t:
            return "&mut " + canon_deployed(t["MutableReference"], addr_map, type_param_names)
        if "Vector" in t:
            return "vector<" + canon_deployed(t["Vector"], addr_map, type_param_names) + ">"
        if "Struct" in t:
            s = t["Struct"]
            addr = addr_map.get(s["address"].lower(), s["address"].lower())
            base = f"{addr}::{s['module']}::{s['name']}"
            args = s.get("typeArguments", [])
            if args:
                base += "<" + ", ".join(
                    canon_deployed(a, addr_map, type_param_names) for a in args
                ) + ">"
            return base
    return repr(t)


def strip_comments(s: str) -> str:
    s = re.sub(r"/\*.*?\*/", "", s, flags=re.DOTALL)
    return re.sub(r"//[^\n]*", "", s)


def find_matching(src: str, start_idx: int, open_ch: str, close_ch: str) -> int:
    depth = 1
    i = start_idx
    while i < len(src) and depth > 0:
        if src[i] == open_ch:
            depth += 1
        elif src[i] == close_ch:
            depth -= 1
        i += 1
    return i


def split_balanced(text: str, sep: str) -> list[str]:
    out, cur, depth = [], "", 0
    for ch in text:
        if ch in "<({[":
            depth += 1
            cur += ch
        elif ch in ">)}]":
            depth -= 1
            cur += ch
        elif ch == sep and depth == 0:
            if cur.strip():
                out.append(cur.strip())
            cur = ""
        else:
            cur += ch
    if cur.strip():
        out.append(cur.strip())
    return out


def find_balanced_end(text: str, start: int, open_ch: str = "<", close_ch: str = ">") -> int:
    if start >= len(text) or text[start] != open_ch:
        return -1
    depth = 1
    i = start + 1
    while i < len(text) and depth > 0:
        if text[i] == open_ch:
            depth += 1
        elif text[i] == close_ch:
            depth -= 1
        i += 1
    return i - 1 if depth == 0 else -1


def parse_type_params(text: str) -> list[tuple[str, list[str], bool]]:
    out = []
    for it in split_balanced(text, ","):
        s = it.strip()
        phantom = s.startswith("phantom ")
        if phantom:
            s = s[len("phantom ") :]
        if ":" in s:
            name, abils = s.split(":", 1)
            ab = sorted(a.strip().lower() for a in abils.split("+") if a.strip())
            out.append((name.strip(), ab, phantom))
        else:
            out.append((s.strip(), [], phantom))
    return out


def parse_use(content: str) -> dict[str, tuple[str, str, str]]:
    aliases: dict[str, tuple[str, str, str]] = {}
    for um in re.finditer(
        r"use\s+([\w:]+)(?:\s*::\s*\{([^}]+)\})?(?:\s+as\s+(\w+))?\s*;", content
    ):
        path, brace, as_alias = um.group(1), um.group(2), um.group(3)
        path_parts = path.split("::")
        addr = ADDR_MAP_LOCAL.get(path_parts[0], "0x" + path_parts[0].upper())
        if brace:
            module_part = "::".join(path_parts[1:]) if len(path_parts) > 1 else ""
            for item in [i.strip() for i in brace.split(",")]:
                am = re.match(r"(\w+)(?:\s+as\s+(\w+))?", item)
                if not am:
                    continue
                orig = am.group(1)
                alias = am.group(2) or orig
                if orig == "Self":
                    last = path_parts[-1]
                    aliases[am.group(2) or last] = (addr, last, last)
                else:
                    aliases[alias] = (addr, module_part or path_parts[-1], orig)
        else:
            if len(path_parts) == 2:
                mod = path_parts[1]
                aliases[as_alias or mod] = (addr, mod, mod)
            elif len(path_parts) >= 3:
                mod = path_parts[-2]
                name = path_parts[-1]
                aliases[as_alias or name] = (addr, mod, name)
    return aliases


def canon_local_type(
    t: str, current_module: str, aliases: dict[str, tuple[str, str, str]], fn_tp_names: list[str]
) -> str:
    t = t.strip()
    if not t:
        return ""
    if t.startswith("&mut "):
        return "&mut " + canon_local_type(t[5:], current_module, aliases, fn_tp_names)
    if t.startswith("&"):
        return "&" + canon_local_type(t[1:].lstrip(), current_module, aliases, fn_tp_names)
    if re.match(r"^vector\s*<", t):
        lt = t.index("<")
        gt = find_balanced_end(t, lt, "<", ">")
        if gt < 0:
            return f"<<UNPARSED:{t[:40]}>>"
        return (
            "vector<"
            + canon_local_type(t[lt + 1 : gt].strip(), current_module, aliases, fn_tp_names)
            + ">"
        )
    if t in ("u8", "u16", "u32", "u64", "u128", "u256", "bool", "address", "signer"):
        return t
    if "<" in t:
        lt = t.index("<")
        gt = find_balanced_end(t, lt, "<", ">")
        if gt < 0:
            return f"<<UNPARSED:{t[:40]}>>"
        base = t[:lt].strip()
        args = split_balanced(t[lt + 1 : gt], ",")
        args_c = [canon_local_type(a, current_module, aliases, fn_tp_names) for a in args]
    else:
        base = t
        args_c = []
    parts = base.split("::")
    if len(parts) == 1:
        name = parts[0]
        if name in fn_tp_names:
            return name + ("<" + ", ".join(args_c) + ">" if args_c else "")
        if name in aliases:
            addr, mod, item = aliases[name]
            base_str = f"{addr}::{mod}::{item}"
        elif name in IMPLICIT_NAME:
            addr, mod, item = IMPLICIT_NAME[name]
            base_str = f"{addr}::{mod}::{item}"
        else:
            base_str = f"0xSELF::{current_module}::{name}"
    else:
        mod_part = parts[0]
        name = parts[-1]
        if mod_part in aliases:
            addr, mod, _ = aliases[mod_part]
            base_str = f"{addr}::{mod}::{name}"
        else:
            base_str = f"0xSELF::{mod_part}::{name}"
    if args_c:
        return base_str + "<" + ", ".join(args_c) + ">"
    return base_str


def parse_module_fn_sigs(path: Path) -> tuple[str | None, dict[str, Any]]:
    content = path.read_text()
    mm = re.search(r"^\s*module\s+ccip::(\w+)", content, re.M)
    if not mm:
        return None, {}
    module_name = mm.group(1)
    src = strip_comments(content)
    aliases = parse_use(content)
    fns: dict[str, Any] = {}
    for fm in re.finditer(r"public(?:\s*\(\s*(\w+)\s*\))?(?:\s+entry)?\s+fun\s+(\w+)", src):
        if fm.group(1) is not None:
            continue
        fname = fm.group(2)
        if fname.endswith("_for_test") or fname == "test_init":
            continue
        start = fm.start()
        if "#[test_only]" in content[max(0, start - 200) : start]:
            continue
        i = fm.end()
        while i < len(src) and src[i] in " \t\n\r":
            i += 1
        tp_list = []
        if i < len(src) and src[i] == "<":
            tp_close = find_matching(src, i + 1, "<", ">")
            tp_list = parse_type_params(src[i + 1 : tp_close - 1])
            i = tp_close
        while i < len(src) and src[i] in " \t\n\r":
            i += 1
        if i >= len(src) or src[i] != "(":
            continue
        p_close = find_matching(src, i + 1, "(", ")")
        params_text = src[i + 1 : p_close - 1]
        tp_names = [n for (n, _, _) in tp_list]
        params = []
        for p in split_balanced(params_text, ","):
            mp = re.match(r"^(?:mut\s+)?(\w+|_)\s*:\s*(.+)$", p, re.DOTALL)
            if mp:
                params.append(
                    (
                        mp.group(1),
                        canon_local_type(mp.group(2).strip(), module_name, aliases, tp_names),
                    )
                )
        i = p_close
        while i < len(src) and src[i] in " \t\n\r":
            i += 1
        ret_canon = ""
        if i < len(src) and src[i] == ":":
            i += 1
            depth = 0
            j = i
            while j < len(src):
                ch = src[j]
                if depth == 0 and ch in "{;":
                    break
                if ch in "<([":
                    depth += 1
                elif ch in ">)]":
                    depth -= 1
                j += 1
            ret_text = src[i:j].strip()
            if ret_text:
                if ret_text.startswith("(") and ret_text.endswith(")"):
                    inner = ret_text[1:-1]
                    parts = split_balanced(inner, ",")
                    parts_c = [
                        canon_local_type(p, module_name, aliases, tp_names) for p in parts
                    ]
                    ret_canon = parts_c[0] if len(parts_c) == 1 else "(" + ", ".join(parts_c) + ")"
                else:
                    ret_canon = canon_local_type(ret_text, module_name, aliases, tp_names)
        fns[fname] = {
            "tp_names": tp_names,
            "tp_constraints": [ab for (_, ab, _) in tp_list],
            "tp_phantom": [ph for (_, _, ph) in tp_list],
            "params": params,
            "return": ret_canon,
        }
    return module_name, fns


def parse_module_structs(path: Path) -> tuple[str | None, dict[str, Any]]:
    content = path.read_text()
    mm = re.search(r"^\s*module\s+ccip::(\w+)", content, re.M)
    if not mm:
        return None, {}
    module_name = mm.group(1)
    src = strip_comments(content)
    aliases = parse_use(content)
    structs: dict[str, Any] = {}
    for sm in re.finditer(r"public\s+struct\s+(\w+)", src):
        sname = sm.group(1)
        i = sm.end()
        while i < len(src) and src[i] in " \t\n\r":
            i += 1
        tp_list = []
        if i < len(src) and src[i] == "<":
            tp_close = find_matching(src, i + 1, "<", ">")
            tp_list = parse_type_params(src[i + 1 : tp_close - 1])
            i = tp_close
        while i < len(src) and src[i] in " \t\n\r":
            i += 1
        if src[i : i + 3] == "has":
            j = i + 3
            while j < len(src) and src[j] != "{":
                j += 1
            abilities = sorted(a.strip() for a in src[i + 3 : j].split(",") if a.strip())
            i = j
        else:
            abilities = []
        if i >= len(src) or src[i] != "{":
            continue
        body_end = find_matching(src, i + 1, "{", "}")
        body = src[i + 1 : body_end - 1]
        tp_names = [n for (n, _, _) in tp_list]
        fields = []
        for line in split_balanced(body, ","):
            mf = re.match(r"^(\w+)\s*:\s*(.+)$", line, re.DOTALL)
            if mf:
                fields.append(
                    (
                        mf.group(1),
                        canon_local_type(mf.group(2).strip(), module_name, aliases, tp_names),
                    )
                )
        structs[sname] = {
            "abilities": abilities,
            "tp_names": tp_names,
            "tp_constraints": [ab for (_, ab, _) in tp_list],
            "tp_phantom": [ph for (_, _, ph) in tp_list],
            "fields": fields,
        }
    return module_name, structs


def parse_local() -> tuple[dict[str, Any], dict[str, Any]]:
    funs: dict[str, Any] = {}
    structs: dict[str, Any] = {}
    for pattern in SRC_GLOB:
        for p in sorted(glob.glob(str(REPO_ROOT / pattern))):
            mname, mod_funs = parse_module_fn_sigs(Path(p))
            if mname:
                funs[mname] = mod_funs
            mname2, mod_structs = parse_module_structs(Path(p))
            if mname2:
                structs[mname2] = mod_structs
    return funs, structs


def struct_abilities(s: dict[str, Any]) -> list[str]:
    ab = s.get("abilities", [])
    if isinstance(ab, dict):
        ab = ab.get("abilities", [])
    return sorted(a.lower() for a in ab)


def struct_tparam_info(s: dict[str, Any]) -> list[tuple[bool, list[str]]]:
    out = []
    for tp in s.get("typeParameters", []):
        c = tp.get("constraints", {})
        if isinstance(c, dict) and "constraints" in c:
            ab = c["constraints"].get("abilities", [])
        elif isinstance(c, dict):
            ab = c.get("abilities", [])
        else:
            ab = []
        out.append((tp.get("isPhantom", False), sorted(a.lower() for a in ab)))
    return out


def fn_tparam_info(fn: dict[str, Any]) -> list[tuple[bool, list[str]]]:
    out = []
    for tp in fn.get("typeParameters", []):
        c = tp.get("constraints", {})
        ab = c.get("abilities", []) if isinstance(c, dict) else []
        out.append((tp.get("isPhantom", False), sorted(a.lower() for a in ab)))
    return out


def is_marker_struct(fields: list[dict[str, Any]]) -> bool:
    return not fields or (len(fields) == 1 and fields[0]["name"] == "dummy_field")


def compare_deployed_to_deployed(
    left: dict[str, Any], right: dict[str, Any], label: str
) -> CompareReport:
    report = CompareReport(label=label)
    left_addr = discover_addr_map(left)
    right_addr = discover_addr_map(right)
    if set(left.keys()) != set(right.keys()):
        report.issues.append(
            f"module set differs: left={sorted(left.keys())} right={sorted(right.keys())}"
        )
    for mod in sorted(set(left.keys()) & set(right.keys())):
        lf, rf = left[mod], right[mod]
        for fname, fd in lf.get("exposedFunctions", {}).items():
            if fd.get("visibility") != "Public":
                continue
            rd = rf.get("exposedFunctions", {}).get(fname)
            if rd is None:
                report.issues.append(f"{mod}::{fname}: missing in right")
                continue
            if rd.get("visibility") != "Public":
                report.issues.append(f"{mod}::{fname}: visibility changed")
            tp = [f"T{i}" for i in range(len(fd.get("typeParameters", [])))]
            d_params = [canon_deployed(p, left_addr, tp) for p in fd.get("parameters", [])]
            r_params = [canon_deployed(p, right_addr, tp) for p in rd.get("parameters", [])]
            d_rets = [canon_deployed(r, left_addr, tp) for r in fd.get("return", [])]
            r_rets = [canon_deployed(r, right_addr, tp) for r in rd.get("return", [])]
            if d_params != r_params or d_rets != r_rets:
                report.issues.append(f"{mod}::{fname}: signature differs")
        for sname, sd in lf.get("structs", {}).items():
            rd = rf.get("structs", {}).get(sname)
            if rd is None:
                report.issues.append(f"{mod}::{sname}: struct missing in right")
                continue
            if struct_abilities(sd) != struct_abilities(rd):
                report.issues.append(f"{mod}::{sname}: abilities differ")
            if is_marker_struct(sd.get("fields", [])) and is_marker_struct(rd.get("fields", [])):
                continue
            tp = [f"T{i}" for i in range(len(sd.get("typeParameters", [])))]
            df = [(f["name"], canon_deployed(f["type"], left_addr, tp)) for f in sd.get("fields", [])]
            rf_ = [(f["name"], canon_deployed(f["type"], right_addr, tp)) for f in rd.get("fields", [])]
            if df != rf_:
                report.issues.append(f"{mod}::{sname}: fields differ")
    report.stats = {
        "modules": len(left),
        "issues": len(report.issues),
    }
    return report


def compare_deployed_to_local(
    deployed: dict[str, Any], local_funs: dict[str, Any], local_structs: dict[str, Any], label: str
) -> CompareReport:
    report = CompareReport(label=label)
    addr_map = discover_addr_map(deployed)
    report.stats["external_addrs"] = {
        k: v for k, v in addr_map.items() if v in ("0xSELF", "0xMCMS")
    }

    pub_fun_count = 0
    marker_struct_count = 0
    checked_struct_count = 0

    for mod in sorted(deployed.keys()):
        if mod not in local_funs or mod not in local_structs:
            report.issues.append(f"MISSING MODULE in local: {mod}")
            continue
        dm = deployed[mod]
        lfns = local_funs[mod]
        lsts = local_structs[mod]

        for sname, sd in dm.get("structs", {}).items():
            ls = lsts.get(sname)
            if ls is None:
                report.issues.append(f"{mod}::{sname}: struct removed/renamed")
                continue
            if struct_abilities(sd) != ls["abilities"]:
                report.issues.append(
                    f"{mod}::{sname}: abilities {struct_abilities(sd)} vs {ls['abilities']}"
                )
            dtps = struct_tparam_info(sd)
            if len(dtps) != len(ls["tp_names"]):
                report.issues.append(f"{mod}::{sname}: tparam count {len(dtps)} vs {len(ls['tp_names'])}")
            else:
                for i, ((d_ph, d_ab), l_ab, l_ph) in enumerate(
                    zip(dtps, ls["tp_constraints"], ls["tp_phantom"])
                ):
                    if set(d_ab) - set(l_ab):
                        report.issues.append(
                            f"{mod}::{sname}: tparam[{i}] constraint tightened deployed={d_ab} local={l_ab}"
                        )
                    if d_ph != l_ph:
                        report.issues.append(f"{mod}::{sname}: tparam[{i}] phantom flag differs")

            dfields = sd.get("fields", [])
            if is_marker_struct(dfields):
                marker_struct_count += 1
                if ls["fields"]:
                    if not (len(ls["fields"]) == 0):
                        report.issues.append(f"{mod}::{sname}: marker struct mismatch")
                continue

            lfields = ls["fields"]
            if len(dfields) != len(lfields):
                report.issues.append(
                    f"{mod}::{sname}: field count {len(dfields)} vs {len(lfields)}"
                )
                continue
            checked_struct_count += 1
            tp_names = ls["tp_names"]
            for (df, (ln, lt)) in zip(dfields, lfields):
                d_canon = canon_deployed(df["type"], addr_map, tp_names)
                if "<<" in lt:
                    report.unparsed.append(f"{mod}::{sname}.{ln}: {lt}")
                    continue
                if df["name"] != ln:
                    report.issues.append(f"{mod}::{sname}: field name {df['name']} vs {ln}")
                if d_canon != lt:
                    report.issues.append(f"{mod}::{sname}.{df['name']}: type {d_canon} vs {lt}")

        for sname in lsts:
            if sname not in dm.get("structs", {}):
                report.additions_structs.append(f"{mod}::{sname}")

        for fname, fd in dm.get("exposedFunctions", {}).items():
            if fd.get("visibility") != "Public":
                continue
            pub_fun_count += 1
            lf = lfns.get(fname)
            if lf is None:
                report.issues.append(f"{mod}::{fname}: function removed/renamed")
                continue
            tp_names = lf["tp_names"]
            dtps = fn_tparam_info(fd)
            if len(dtps) != len(lf["tp_names"]):
                report.issues.append(f"{mod}::{fname}: tparam count")
            else:
                for i, ((d_ph, d_ab), l_ab, l_ph) in enumerate(
                    zip(dtps, lf["tp_constraints"], lf["tp_phantom"])
                ):
                    if set(d_ab) - set(l_ab):
                        report.issues.append(
                            f"{mod}::{fname}: tparam[{i}] constraint tightened deployed={d_ab} local={l_ab}"
                        )
                    if d_ph != l_ph:
                        report.issues.append(f"{mod}::{fname}: tparam[{i}] phantom differs")

            d_params = [canon_deployed(p, addr_map, tp_names) for p in fd.get("parameters", [])]
            d_ret_list = [canon_deployed(r, addr_map, tp_names) for r in fd.get("return", [])]
            if not d_ret_list:
                d_return = ""
            elif len(d_ret_list) == 1:
                d_return = d_ret_list[0]
            else:
                d_return = "(" + ", ".join(d_ret_list) + ")"

            if len(d_params) != len(lf["params"]):
                report.issues.append(
                    f"{mod}::{fname}: param count {len(d_params)} vs {len(lf['params'])}"
                )
                continue
            for idx, (d_t, (pname, l_t)) in enumerate(zip(d_params, lf["params"])):
                if "<<" in l_t:
                    report.unparsed.append(f"{mod}::{fname} arg[{idx}] ({pname}): {l_t}")
                    continue
                if d_t != l_t:
                    report.issues.append(f"{mod}::{fname} arg[{idx}] ({pname}): {d_t} vs {l_t}")
            if "<<" in lf["return"]:
                report.unparsed.append(f"{mod}::{fname} return: {lf['return']}")
            elif d_return != lf["return"]:
                report.issues.append(f"{mod}::{fname} return: {d_return!r} vs {lf['return']!r}")

        deployed_public = {
            fname
            for fname, fd in dm.get("exposedFunctions", {}).items()
            if fd.get("visibility") == "Public"
        }
        for fname in lfns:
            if fname not in deployed_public:
                report.additions_funs.append(f"{mod}::{fname}")

    report.stats.update(
        {
            "public_functions": pub_fun_count,
            "non_marker_structs_checked": checked_struct_count,
            "marker_structs": marker_struct_count,
            "issues": len(report.issues),
            "unparsed": len(report.unparsed),
            "additive_funs": len(report.additions_funs),
            "additive_structs": len(report.additions_structs),
        }
    )
    return report


def load_json(path: Path) -> dict[str, Any]:
    data = json.loads(path.read_text())
    if "result" in data:
        return data["result"]
    return data


def save_json(path: Path, obj: Any) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(json.dumps(obj, indent=2))


def remap_package_addresses(modules: dict[str, Any], replacements: dict[str, str]) -> dict[str, Any]:
    text = json.dumps(modules)
    for old, new in replacements.items():
        text = text.replace(old, new)
        text = text.replace(old.upper(), new)
    return json.loads(text)


def main() -> int:
    parser = argparse.ArgumentParser(description="Compare CCIP upgrade compatibility")
    parser.add_argument("--network", choices=["mainnet", "testnet"], default="mainnet")
    parser.add_argument("--v1-package", required=True)
    parser.add_argument("--latest-package", default="")
    parser.add_argument("--upgrade-cap", default="")
    parser.add_argument("--v1-json", type=Path, default=None)
    parser.add_argument("--latest-json", type=Path, default=None)
    parser.add_argument("--out-dir", type=Path, default=REPO_ROOT / ".upgrade-analysis")
    parser.add_argument("--fetch", action="store_true")
    parser.add_argument("--json-report", type=Path, default=None)
    args = parser.parse_args()

    out_dir = args.out_dir
    v1_path = args.v1_json or out_dir / f"ccip_{args.network}_v1_modules.json"
    latest_path = args.latest_json or out_dir / f"ccip_{args.network}_latest_modules.json"

    upgrade_meta = {}
    if args.fetch:
        try:
            v1_modules = fetch_normalized_modules(args.v1_package, args.network)
            save_json(v1_path, {"jsonrpc": "2.0", "id": 1, "result": v1_modules})
            print(f"Fetched v1 modules -> {v1_path} ({len(v1_modules)} modules)")
        except RuntimeError as exc:
            print(f"ERROR: {exc}", file=sys.stderr)
            return 1
        if args.upgrade_cap:
            upgrade_meta = fetch_upgrade_cap(args.v1_package, args.upgrade_cap, args.network)
            save_json(out_dir / f"ccip_{args.network}_upgrade_cap.json", upgrade_meta)
            print(f"UpgradeCap version={upgrade_meta.get('version')} policy={upgrade_meta.get('policy')}")

    if not v1_path.exists():
        print(f"Missing v1 JSON: {v1_path}", file=sys.stderr)
        return 1

    v1 = load_json(v1_path)
    latest = v1
    latest_note = "latest same as v1 (single published version)"
    if args.latest_package and args.latest_package != args.v1_package:
        if args.fetch:
            try:
                latest = fetch_normalized_modules(args.latest_package, args.network)
                save_json(latest_path, {"jsonrpc": "2.0", "id": 1, "result": latest})
            except RuntimeError as exc:
                print(f"WARN: could not fetch latest package: {exc}", file=sys.stderr)
        elif latest_path.exists():
            latest = load_json(latest_path)
            latest_note = f"loaded {latest_path}"
        else:
            print("WARN: latest package specified but no JSON available", file=sys.stderr)

    local_funs, local_structs = parse_local()
    v1_vs_latest = compare_deployed_to_deployed(v1, latest, "v1_vs_latest")
    v1_vs_local = compare_deployed_to_local(v1, local_funs, local_structs, "v1_vs_local")

    report = {
        "network": args.network,
        "v1_package": args.v1_package,
        "latest_package": args.latest_package or args.v1_package,
        "upgrade_cap": upgrade_meta or None,
        "latest_note": latest_note,
        "v1_vs_latest": {
            "issues": v1_vs_latest.issues,
            "stats": v1_vs_latest.stats,
        },
        "v1_vs_local": {
            "issues": v1_vs_local.issues,
            "unparsed": v1_vs_local.unparsed,
            "additions_funs": v1_vs_local.additions_funs,
            "additions_structs": v1_vs_local.additions_structs,
            "stats": v1_vs_local.stats,
        },
    }

    if args.json_report:
        save_json(args.json_report, report)
    else:
        save_json(out_dir / f"ccip_{args.network}_compat_report.json", report)

    print(json.dumps(report, indent=2))
    ok = not v1_vs_latest.issues and not v1_vs_local.issues and not v1_vs_local.unparsed
    return 0 if ok else 2


if __name__ == "__main__":
    raise SystemExit(main())
