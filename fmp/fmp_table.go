package fmp

func (ctx *FmpFile) Tables() []*FmpTable {
	tables := make([]*FmpTable, 0)

	for key, ent := range *ctx.Dictionary {
		if key != 3 {
			continue
		}
		debug("Found a 3")

		for key, ent = range *ent.Children {
			if key != 16 {
				continue
			}
			debug("Found a 3.16")

			for key, ent = range *ent.Children {
				if key != 5 {
					continue
				}
				debug("Found a 3.16.5")

				for tablePath := range *ent.Children {
					if key >= 128 {
						continue
					}

					// Found a table!
					debug("Found a table at 3.16.5.%d", tablePath)
				}
			}
		}
	}

	// metaDict := ctx.Dictionary.get([]uint64{4, 1, 7})
	// if metaDict == nil {
	// 	return tables
	// }
	// for _, meta := range *metaDict.Children {
	// 	name := decodeByteSeq(meta.Children.get([]uint64{16}).Value)
	// 	table := &FmpTable{Name: name}
	// 	tables = append(tables, table)
	// }

	return tables
}
