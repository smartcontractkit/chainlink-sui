package chainwriter

import (
	"fmt"

	"github.com/pattonkan/sui-go/sui"

	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
)

func (p *PTBConstructor) ResolveGenericTypeTags(
	params []codec.SuiFunctionParam,
	arguments Arguments,
) ([]sui.TypeTag, error) {
	if len(params) == 0 {
		return nil, nil
	}

	genericParams := []codec.SuiFunctionParam{}
	for _, param := range params {
		if param.IsGeneric {
			genericParams = append(genericParams, param)
		}
	}

	// Use a map to track unique type tags
	seen := make(map[string]struct{})
	out := make([]sui.TypeTag, 0, len(genericParams))

	p.log.Debugw("Resolving generic type tags", "genericParams", genericParams, "arguments", arguments)

	for _, g := range genericParams {
		if g.Name == "" {
			return nil, fmt.Errorf("generic argument must specify a parameter name")
		}

		typeTag, ok := arguments.ArgTypes[g.Name]
		if !ok {
			return nil, fmt.Errorf("generic arg refers to parameter %q that does not exist", g.Name)
		}

		// Skip if we've already seen this type tag
		if _, exists := seen[typeTag]; exists {
			continue
		}

		tag, err := sui.NewTypeTag(typeTag)
		if err != nil {
			return nil, fmt.Errorf("invalid type tag for parameter %q: %w", g.Name, err)
		}

		seen[typeTag] = struct{}{}
		out = append(out, *tag)
	}

	return out, nil
}
