# Textbench

## Test Data

The `/data` directory contains text files used for benchmarking and testing the
text evaluation functionality.

| File              | Description                                                                                                          |
|-------------------|----------------------------------------------------------------------------------------------------------------------|
| `moby-dick-0.txt` | Full text of Herman Melville's "Moby-Dick" novel, used as reference corpus for text comparison and evaluation tests. |
| `moby-dick-1.txt` |                                                                                                                      |

The test data files are referenced in the main application (see `cmd/main.go`)
to demonstrate and benchmark the text normalization and comparison pipeline.
