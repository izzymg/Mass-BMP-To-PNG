# Mass-BMP-To-PNG

Converts all BMPs in the given directory to PNGs

Options:

* `--input="path"` Directory to process bmps in. Defaults to current directory.
* `--output="path"` Directory to output pngs to. Defaults to input directory.
* `--clean` Remove BMPs after processing.
* `--c=5` Sets number of concurrent operations to 5.
* `--silent` Prevents stdout.
* `--help`

## Usage

You can download the compiled binaries in the [releases](https://github.com/izzymg/Mass-BMP-To-PNG/releases) section

For windows, you can right click the exe, create a shortcut, open its properties add the run flags to the "Target" section

e.g `C:\Users\You\Downloads\Mass-BMP-To-PNG.exe --input="D:\Games\SteamLibrary\steamapps\common\Skyrim" --output="D:\Pictures\SkyrimScreenshots"`
