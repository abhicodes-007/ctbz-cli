<#
.SYNOPSIS
  Define as variáveis de ambiente do ctbz para a sessão atual do PowerShell.

.EXAMPLE
  # IMPORTANTE: use dot-sourcing (ponto + espaço) para as variáveis valerem na sua sessão
  . .\scripts\ctbz-env.ps1
  . .\scripts\ctbz-env.ps1 -User voce@email.com -Cnpj 12345678000190 -Output json
  . .\scripts\ctbz-env.ps1 -OtpCmd "$PWD/scripts/otp-gmail-gws.sh"
#>
param(
    [string]$User     = $env:CTBZ_USER,
    [string]$Cnpj     = $env:CTBZ_CNPJ,
    [string]$OtpCmd   = $env:CTBZ_OTP_CMD,
    [string]$OtpTimeout = "3m",
    [ValidateSet("table", "json", "csv")]
    [string]$Output   = "table",
    [string]$CtbzHome = $env:CTBZ_HOME,
    [switch]$VerboseOtp
)

# Login: pergunta se não foi informado
if (-not $User) { $User = Read-Host "CTBZ_USER (e-mail ou CPF)" }
$env:CTBZ_USER = $User

# Senha: lida sem eco e mantida só na memória desta sessão (nunca gravada em disco)
if (-not $env:CTBZ_PASSWORD) {
    $secure = Read-Host "CTBZ_PASSWORD" -AsSecureString
    $bstr = [Runtime.InteropServices.Marshal]::SecureStringToBSTR($secure)
    try { $env:CTBZ_PASSWORD = [Runtime.InteropServices.Marshal]::PtrToStringBSTR($bstr) }
    finally { [Runtime.InteropServices.Marshal]::ZeroFreeBSTR($bstr) }
}

# Opcionais: só define se tiver valor
if ($Cnpj)     { $env:CTBZ_CNPJ = $Cnpj }
if ($CtbzHome) { $env:CTBZ_HOME = $CtbzHome }

# OTP automático. A CLI executa `sh -c "<comando>"`, então no Windows
# precisa de `sh` no PATH (Git for Windows ou WSL).
if ($OtpCmd) {
    if (-not (Get-Command sh -ErrorAction SilentlyContinue)) {
        Write-Warning "CTBZ_OTP_CMD usa 'sh', que não está no PATH. Instale o Git for Windows e adicione '<git>\usr\bin' ao PATH."
    }
    $env:CTBZ_OTP_CMD = $OtpCmd
}
$env:CTBZ_OTP_TIMEOUT = $OtpTimeout
$env:CTBZ_OUTPUT      = $Output
if ($VerboseOtp) { $env:CTBZ_VERBOSE = "1" } else { Remove-Item Env:CTBZ_VERBOSE -ErrorAction SilentlyContinue }

# Garante que o ctbz.exe da raiz do repo esteja no PATH desta sessão
$repoRoot = Split-Path -Parent $PSScriptRoot
if (Test-Path (Join-Path $repoRoot "ctbz.exe")) {
    if (($env:Path -split ';') -notcontains $repoRoot) { $env:Path += ";$repoRoot" }
}

Write-Host "Variáveis do ctbz definidas nesta sessão:" -ForegroundColor Green
Get-ChildItem Env:CTBZ_* |
    Where-Object { $_.Name -ne "CTBZ_PASSWORD" } |
    Sort-Object Name | Format-Table Name, Value -AutoSize
Write-Host "CTBZ_PASSWORD  = ******** (oculta)"
Write-Host "Teste com: ctbz login" -ForegroundColor Cyan
