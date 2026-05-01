# Textbench

## Test Data

The `/data` directory contains text files used for benchmarking and testing the
text evaluation functionality.

Every file in the Moby Dick directory with a number on the end, represents how many mistakes or likely issues there are in the file compared to defy with the lowest number for the file names. A higher post fixed number represents a higher error rate.
                   |

The test data files are referenced in the main application (see `cmd/main.go`)
to demonstrate and benchmark the text normalization and comparison pipeline.
