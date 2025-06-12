package fmp

func (ctx *FmpFile) Tables() []*FmpTable {
	tables := make([]*FmpTable, 0)

	for _, chunk := range ctx.Chunks {
		if chunk.Key != 3 || chunk.Type != FMP_CHUNK_SIMPLE_KEY_VALUE {
			continue
		}

		for _, chunk = range ctx.Chunks {
			if chunk.Key != 16 {
				continue
			}

			for _, chunk = range ctx.Chunks {
				if chunk.Key != 5 {
					continue
				}

				for tablePath, chunk := range ctx.Chunks {
					if chunk.Key >= 128 {
						continue
					}

					// Found a table!
					println("Found one at", tablePath)
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
