$Targets = @("linux/amd64", "windows/amd64")
$OutputDir = ".\build"

New-Item -ItemType Directory -Force -Path $OutputDir

foreach ($Target in $Targets) {
    $Parts = $Target -split "/"
    $OS = $Parts[0]
    $Arch = $Parts[1]
    $OSArch = "$OS_$Arch"
    $Extension = ""

    if ($OS -eq "windows") {
        $Extension = ".exe"
    }

    Write-Host "Building $OS/$Arch..."

    Push-Location ./compare/src
    Write-Host "Building compare..."
    $Env:GOOS = $OS
    $Env:GOARCH = $Arch
    go build -o "../../$OutputDir/$OSArch/Compare$Extension"
    Pop-Location

    Push-Location ./filter/src
    Write-Host "Building filter..."
    $Env:GOOS = $OS
    $Env:GOARCH = $Arch
    go build -o "../../$OutputDir/$OSArch/Filter$Extension"
    Pop-Location

    Push-Location ./browser/src
    Write-Host "Building browser..."
    $Env:GOOS = $OS
    $Env:GOARCH = $Arch
    go build -o "../../$OutputDir/$OSArch/Browser$Extension"
    Pop-Location
}

Push-Location $OutputDir
Get-ChildItem -Directory | ForEach-Object {
    $Folder = $_.FullName
    Compress-Archive -Path "$Folder\*" -DestinationPath "$Folder.zip" -Force
}
Pop-Location
