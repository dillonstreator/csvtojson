package csvtojson

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

// Copy copies src csv content to dst as json objects after applying the specified optional transformations
func Copy(ctx context.Context, src io.Reader, dst io.Writer, options ...option) error {
	cfg := newConfig(options...)

	csvReader := csv.NewReader(src)
	csvReader.Comma = cfg.separator

	headers, err := csvReader.Read()
	if err != nil {
		return err
	}

	err = validateConfig(cfg, headers)
	if err != nil {
		return err
	}

	jsonEncoder := json.NewEncoder(dst)

readLoop:
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		default:
			record, err := csvReader.Read()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break readLoop
				}

				return err
			}

			out := map[string]string{}

			for idx, fieldValue := range record {
				var key string

				if cfg.fields != nil {
					var ok bool

					key, ok = cfg.fields[headers[idx]]
					if !ok {
						continue
					}
				} else {
					key = headers[idx]
				}

				out[key] = fieldValue

				if cfg.inject != nil {
					for k, v := range cfg.inject {
						out[k] = v
					}
				}
			}

			err = jsonEncoder.Encode(out)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func validateConfig(cfg *config, headers []string) error {
	if cfg.fields != nil {
		err := validateFields(cfg.fields, headers)
		if err != nil {
			return err
		}
	}

	if cfg.inject != nil {
		err := validateInject(cfg.inject, cfg.fields, headers)
		if err != nil {
			return err
		}
	}

	return nil
}

func validateFields(fields map[string]string, headers []string) error {
	missingFields := []string{}

	for inKey := range fields {
		var found bool

		for _, header := range headers {
			if header == inKey {
				found = true
				break
			}
		}

		if !found {
			missingFields = append(missingFields, inKey)
		}
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("missing header(s) in incoming csv: %s", strings.Join(missingFields, ", "))
	}

	return nil
}

func validateInject(inject, fields map[string]string, headers []string) error {
	clashingFields := []string{}

	for key := range inject {
		var clash bool

		for _, header := range headers {
			if fields != nil {
				if fields[header] == key {
					clash = true
					break
				}
			} else {
				if header == key {
					clash = true
					break
				}
			}
		}

		if clash {
			clashingFields = append(clashingFields, key)
		}
	}

	if len(clashingFields) > 0 {
		return fmt.Errorf("clashing key(s) in injection: %s", strings.Join(clashingFields, ", "))
	}

	return nil
}
