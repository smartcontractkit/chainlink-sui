package chainaccessor

import (
	"context"
	"fmt"
	"slices"
	"sort"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"

	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/database"
)

// defaultEventQueryLimit bounds how many events we pull from the store before
// filtering in memory. The store also enforces its own MaxEventsQueryLimit.
const defaultEventQueryLimit = 1000

// queryEvents reads events of the given type for a bound contract from the
// indexer-populated event store. Events are returned in store order; callers
// filter/sort as needed. Confidence is accepted for interface parity but the
// store only contains indexed (effectively finalized) events.
func (a *SuiAccessor) queryEvents(
	ctx context.Context,
	contractName, module, eventType string,
	expressions []query.Expression,
	limit uint64,
) ([]database.EventRecord, error) {
	packageAddr, err := a.bindings.getPackageAddress(contractName)
	if err != nil {
		return nil, err
	}

	// Event handle mirrors the ChainReader: "<package>::<module>::<EventType>".
	eventHandle := fmt.Sprintf("%s::%s::%s", packageAddr, module, eventType)

	limitAndSort := query.LimitAndSort{Limit: query.Limit{Count: limit}}
	records, err := a.dbStore.QueryEvents(ctx, packageAddr, eventHandle, expressions, limitAndSort)
	if err != nil {
		return nil, fmt.Errorf("querying %s events: %w", eventType, err)
	}
	return records, nil
}

// --- SourceAccessor (event reads) ---

// MsgsBetweenSeqNums returns all messages sent to the provided destination chain
// between the provided sequence numbers, sorted ascending by sequence number.
//
// Event(onramp::CCIPMessageSent) on OnRamp.
func (a *SuiAccessor) MsgsBetweenSeqNums(ctx context.Context, dest ccipocr3.ChainSelector, seqNumRange ccipocr3.SeqNumRange) ([]ccipocr3.Message, error) {
	records, err := a.queryEvents(ctx, ContractNameOnRamp, "onramp", EventNameCCIPMessageSent, nil, defaultEventQueryLimit)
	if err != nil {
		return nil, err
	}

	msgs := make([]ccipocr3.Message, 0, len(records))
	for i := range records {
		msg, ok, err := decodeCCIPMessageSent(records[i].Data)
		if err != nil {
			a.lggr.Warnw("skipping undecodable CCIPMessageSent event", "err", err, "txDigest", records[i].TxDigest)
			continue
		}
		if !ok {
			continue
		}
		if msg.Header.DestChainSelector != dest {
			continue
		}
		if !seqNumRange.Contains(msg.Header.SequenceNumber) {
			continue
		}
		msgs = append(msgs, msg)
	}

	sort.Slice(msgs, func(i, j int) bool {
		return msgs[i].Header.SequenceNumber < msgs[j].Header.SequenceNumber
	})
	return msgs, nil
}

// LatestMessageTo returns the sequence number of the most recent message sent to
// the given destination.
//
// Event(onramp::CCIPMessageSent) on OnRamp.
func (a *SuiAccessor) LatestMessageTo(ctx context.Context, dest ccipocr3.ChainSelector) (ccipocr3.SeqNum, error) {
	records, err := a.queryEvents(ctx, ContractNameOnRamp, "onramp", EventNameCCIPMessageSent, nil, defaultEventQueryLimit)
	if err != nil {
		return 0, err
	}

	var latest ccipocr3.SeqNum
	var found bool
	for i := range records {
		msg, ok, err := decodeCCIPMessageSent(records[i].Data)
		if err != nil || !ok {
			continue
		}
		if msg.Header.DestChainSelector != dest {
			continue
		}
		if !found || msg.Header.SequenceNumber > latest {
			latest = msg.Header.SequenceNumber
			found = true
		}
	}

	if !found {
		return 0, fmt.Errorf("no CCIPMessageSent events found for dest chain %d", dest)
	}
	return latest, nil
}

// --- DestinationAccessor (event reads) ---

// CommitReportsGTETimestamp reads CommitReportAccepted events at or after the
// given timestamp, up to the provided limit.
//
// Event(offramp::CommitReportAccepted) on OffRamp.
func (a *SuiAccessor) CommitReportsGTETimestamp(
	ctx context.Context,
	ts time.Time,
	_ primitives.ConfidenceLevel,
	limit int,
) ([]ccipocr3.CommitPluginReportWithMeta, error) {
	records, err := a.queryEvents(ctx, ContractNameOffRamp, "offramp", EventNameCommitReportAccepted, nil, defaultEventQueryLimit)
	if err != nil {
		return nil, err
	}

	reports := make([]ccipocr3.CommitPluginReportWithMeta, 0, len(records))
	for i := range records {
		rec := records[i]
		//nolint:gosec // block timestamps are non-negative
		recTs := time.Unix(int64(rec.BlockTimestamp), 0)
		if recTs.Before(ts) {
			continue
		}

		report, err := decodeCommitReport(rec.Data)
		if err != nil {
			a.lggr.Warnw("skipping undecodable CommitReportAccepted event", "err", err, "txDigest", rec.TxDigest)
			continue
		}

		reports = append(reports, ccipocr3.CommitPluginReportWithMeta{
			Report:    report,
			Timestamp: recTs,
			BlockNum:  rec.BlockVersion,
		})
	}

	sort.Slice(reports, func(i, j int) bool {
		return reports[i].Timestamp.Before(reports[j].Timestamp)
	})

	if limit > 0 && len(reports) > limit {
		reports = reports[:limit]
	}
	return reports, nil
}

// ExecutedMessages returns, for each requested source chain, the sequence numbers
// that have an ExecutionStateChanged event within the requested ranges.
//
// Event(offramp::ExecutionStateChanged) on OffRamp.
func (a *SuiAccessor) ExecutedMessages(
	ctx context.Context,
	ranges map[ccipocr3.ChainSelector][]ccipocr3.SeqNumRange,
	_ primitives.ConfidenceLevel,
) (map[ccipocr3.ChainSelector][]ccipocr3.SeqNum, error) {
	records, err := a.queryEvents(ctx, ContractNameOffRamp, "offramp", EventNameExecutionStateChanged, nil, defaultEventQueryLimit)
	if err != nil {
		return nil, err
	}

	// Index executed (source, seqNum) pairs from the events.
	executed := make(map[ccipocr3.ChainSelector]map[ccipocr3.SeqNum]struct{})
	for i := range records {
		src, seq, err := decodeExecutionStateChanged(records[i].Data)
		if err != nil {
			a.lggr.Warnw("skipping undecodable ExecutionStateChanged event", "err", err, "txDigest", records[i].TxDigest)
			continue
		}
		if _, ok := executed[src]; !ok {
			executed[src] = make(map[ccipocr3.SeqNum]struct{})
		}
		executed[src][seq] = struct{}{}
	}

	out := make(map[ccipocr3.ChainSelector][]ccipocr3.SeqNum, len(ranges))
	for src, srcRanges := range ranges {
		seqSet := executed[src]
		if len(seqSet) == 0 {
			continue
		}
		var seqs []ccipocr3.SeqNum
		for seq := range seqSet {
			if seq.IsWithinRanges(srcRanges) {
				seqs = append(seqs, seq)
			}
		}
		if len(seqs) == 0 {
			continue
		}
		slices.Sort(seqs)
		out[src] = seqs
	}
	return out, nil
}

// --- decoders ---

// decodeCCIPMessageSent decodes a CCIPMessageSent event into a ccipocr3.Message.
// The Header (the fields used for filtering and ordering) is fully populated;
// the remaining scalar fields are best-effort. Token-amount mapping is left for
// follow-up work (see plan).
func decodeCCIPMessageSent(data map[string]any) (ccipocr3.Message, bool, error) {
	msgMap, err := mapField(data, "message")
	if err != nil {
		return ccipocr3.Message{}, false, err
	}
	headerMap, err := mapField(msgMap, "header")
	if err != nil {
		return ccipocr3.Message{}, false, err
	}

	srcSel, err := uint64Field(headerMap, "source_chain_selector")
	if err != nil {
		return ccipocr3.Message{}, false, err
	}
	destSel, err := uint64Field(headerMap, "dest_chain_selector")
	if err != nil {
		return ccipocr3.Message{}, false, err
	}
	seqNum, err := uint64Field(headerMap, "sequence_number")
	if err != nil {
		return ccipocr3.Message{}, false, err
	}
	nonce, err := uint64Field(headerMap, "nonce")
	if err != nil {
		return ccipocr3.Message{}, false, err
	}
	messageID, err := bytesField(headerMap, "message_id")
	if err != nil {
		return ccipocr3.Message{}, false, err
	}

	header := ccipocr3.RampMessageHeader{
		MessageID:           toBytes32(messageID),
		SourceChainSelector: ccipocr3.ChainSelector(srcSel),
		DestChainSelector:   ccipocr3.ChainSelector(destSel),
		SequenceNumber:      ccipocr3.SeqNum(seqNum),
		Nonce:               nonce,
	}

	msg := ccipocr3.Message{Header: header}

	// Best-effort scalar fields; absence is not fatal.
	if sender, err := bytesField(msgMap, "sender"); err == nil {
		msg.Sender = sender
	}
	if receiver, err := bytesField(msgMap, "receiver"); err == nil {
		msg.Receiver = receiver
	}
	if payload, err := bytesField(msgMap, "data"); err == nil {
		msg.Data = payload
	}
	if extraArgs, err := bytesField(msgMap, "extra_args"); err == nil {
		msg.ExtraArgs = extraArgs
	}
	if feeToken, err := bytesField(msgMap, "fee_token"); err == nil {
		msg.FeeToken = feeToken
	}
	if rawFee, err := field(msgMap, "fee_token_amount"); err == nil {
		if amt, err := asBigInt(rawFee); err == nil {
			msg.FeeTokenAmount = ccipocr3.NewBigInt(amt)
		}
	}
	if rawJuels, err := field(msgMap, "fee_value_juels"); err == nil {
		if amt, err := asBigInt(rawJuels); err == nil {
			msg.FeeValueJuels = ccipocr3.NewBigInt(amt)
		}
	}

	return msg, true, nil
}

// decodeExecutionStateChanged extracts (source chain selector, sequence number)
// from an ExecutionStateChanged event.
func decodeExecutionStateChanged(data map[string]any) (ccipocr3.ChainSelector, ccipocr3.SeqNum, error) {
	src, err := uint64Field(data, "source_chain_selector")
	if err != nil {
		return 0, 0, err
	}
	seq, err := uint64Field(data, "sequence_number")
	if err != nil {
		return 0, 0, err
	}
	return ccipocr3.ChainSelector(src), ccipocr3.SeqNum(seq), nil
}

// decodeCommitReport decodes a CommitReportAccepted event into a
// ccipocr3.CommitPluginReport. Merkle roots are mapped; price updates and RMN
// signatures are left for follow-up work (see plan).
func decodeCommitReport(data map[string]any) (ccipocr3.CommitPluginReport, error) {
	report := ccipocr3.CommitPluginReport{}

	blessed, err := decodeMerkleRoots(data, "blessed_merkle_roots")
	if err != nil {
		return ccipocr3.CommitPluginReport{}, err
	}
	unblessed, err := decodeMerkleRoots(data, "unblessed_merkle_roots")
	if err != nil {
		return ccipocr3.CommitPluginReport{}, err
	}
	report.BlessedMerkleRoots = blessed
	report.UnblessedMerkleRoots = unblessed
	return report, nil
}

func decodeMerkleRoots(data map[string]any, key string) ([]ccipocr3.MerkleRootChain, error) {
	raw, err := field(data, key)
	if err != nil {
		// Treat absence as empty rather than fatal.
		//nolint:nilerr // missing roots is a valid (price-only) report
		return nil, nil
	}
	list, ok := raw.([]any)
	if !ok {
		return nil, fmt.Errorf("field %q is not a list", key)
	}

	roots := make([]ccipocr3.MerkleRootChain, 0, len(list))
	for _, item := range list {
		m, err := asMap(item)
		if err != nil {
			return nil, err
		}
		srcSel, err := uint64Field(m, "source_chain_selector")
		if err != nil {
			return nil, err
		}
		minSeq, err := uint64Field(m, "min_seq_nr")
		if err != nil {
			return nil, err
		}
		maxSeq, err := uint64Field(m, "max_seq_nr")
		if err != nil {
			return nil, err
		}
		onRamp, err := bytesField(m, "on_ramp_address")
		if err != nil {
			return nil, err
		}
		root, err := bytesField(m, "merkle_root")
		if err != nil {
			return nil, err
		}
		roots = append(roots, ccipocr3.MerkleRootChain{
			ChainSel:      ccipocr3.ChainSelector(srcSel),
			OnRampAddress: onRamp,
			SeqNumsRange:  ccipocr3.NewSeqNumRange(ccipocr3.SeqNum(minSeq), ccipocr3.SeqNum(maxSeq)),
			MerkleRoot:    toBytes32(root),
		})
	}
	return roots, nil
}
