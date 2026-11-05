# install-skill.ps1 - install this repo as the DSH skill "nsgo".
#
#   pwsh -File tools\install-skill.ps1
#
# The repo root is itself a valid skill directory (SKILL.md + references/), so
# installing means pointing ~/.dsh/skills/nsgo at the checkout: a directory
# junction when Windows allows it, a copy of the two paths otherwise.

$ErrorActionPreference = "Stop"

$repo = Split-Path -Parent $PSScriptRoot
$skills = Join-Path $env:USERPROFILE ".dsh\skills"
$link = Join-Path $skills "nsgo"

if (-not (Test-Path (Join-Path $repo "SKILL.md"))) {
    throw "SKILL.md not found next to $PSScriptRoot - run this from a checkout of the repo."
}

New-Item -ItemType Directory -Path $skills -Force | Out-Null

function Test-Locked([string]$path) {
    # A directory that is being deleted (pending) answers "access denied" to
    # everything, including reads of its own children.
    if (-not (Test-Path $path)) { return $false }
    try {
        Get-ChildItem -LiteralPath $path -Force -ErrorAction Stop | Out-Null
        return $false
    } catch {
        return $true
    }
}

if (Test-Path $link) {
    $item = Get-Item -LiteralPath $link -Force
    if ($item.LinkType -and $item.Target -contains $repo) {
        "already installed: $link -> $repo"
        exit 0
    }
    if (Test-Locked $link) {
        "WARNING: $link is locked by a running process (a delete is pending)."
        "Close the app that watches ~/.dsh/skills and run this script again."
        exit 1
    }
    Remove-Item -LiteralPath $link -Recurse -Force
}

$made = $false
try {
    cmd /c mklink /J "$link" "$repo" | Out-Null
    $made = Test-Path (Join-Path $link "SKILL.md")
} catch {
    $made = $false
}

if (-not $made) {
    New-Item -ItemType Directory -Path $link -Force | Out-Null
    Copy-Item -LiteralPath (Join-Path $repo "SKILL.md") -Destination $link -Force
    Copy-Item -LiteralPath (Join-Path $repo "references") -Destination $link -Recurse -Force
    "installed by copy: $link (re-run this script after editing SKILL.md)"
} else {
    "installed by junction: $link -> $repo"
}

Get-ChildItem -LiteralPath $link -Force | ForEach-Object { "  " + $_.Name }
