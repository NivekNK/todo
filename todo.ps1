$goLangId = winget list --id GoLang.Go
$goLangInstalledText = $goLangId -split '\r?\n' | Select-Object -Last 1

if ($goLangInstalledText -eq "No installed package found matching input criteria.") {
    Write-Output "GoLang.Go not found. Installing..."
    winget install GoLang.Go
    $Env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")       
    Write-Output "GoLang.Go installed."
}

go version
