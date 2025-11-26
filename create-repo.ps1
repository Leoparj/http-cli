# Script to create GitHub repository via API
param(
    [string]$Token = $env:GITHUB_TOKEN,
    [string]$Username = "neore",
    [string]$RepoName = "http-cli"
)

$description = "A colorful, lightweight command-line HTTP client similar to Postman. Built with Go for fast performance and beautiful terminal output."
$topics = @("go", "cli", "http-client", "rest-api", "postman-alternative", "terminal-tool", "http-requests", "api-testing", "command-line", "developer-tools")

if (-not $Token) {
    Write-Host "GitHub token not found. Please set GITHUB_TOKEN environment variable or provide it as parameter." -ForegroundColor Yellow
    Write-Host ""
    Write-Host "To create the repository manually:" -ForegroundColor Cyan
    Write-Host "1. Go to https://github.com/new" -ForegroundColor White
    Write-Host "2. Repository name: $RepoName" -ForegroundColor White
    Write-Host "3. Description: $description" -ForegroundColor White
    Write-Host "4. Set to Public" -ForegroundColor White
    Write-Host "5. DO NOT initialize with README, .gitignore, or license" -ForegroundColor White
    Write-Host "6. Click 'Create repository'" -ForegroundColor White
    Write-Host ""
    Write-Host "Then run:" -ForegroundColor Yellow
    Write-Host "  git remote add origin https://github.com/$Username/$RepoName.git" -ForegroundColor Green
    Write-Host "  git push -u origin main" -ForegroundColor Green
    exit
}

$headers = @{
    "Authorization" = "token $Token"
    "Accept" = "application/vnd.github.v3+json"
}

$body = @{
    name = $RepoName
    description = $description
    private = $false
    auto_init = $false
} | ConvertTo-Json

Write-Host "Creating repository $RepoName on GitHub..." -ForegroundColor Cyan

try {
    $response = Invoke-RestMethod -Uri "https://api.github.com/user/repos" -Method Post -Headers $headers -Body $body -ContentType "application/json"
    Write-Host "Repository created successfully!" -ForegroundColor Green
    Write-Host "URL: $($response.html_url)" -ForegroundColor Cyan
    
    # Add topics
    Write-Host "Adding topics..." -ForegroundColor Cyan
    $topicsBody = @{
        names = $topics
    } | ConvertTo-Json
    
    try {
        Invoke-RestMethod -Uri "https://api.github.com/repos/$Username/$RepoName/topics" -Method Put -Headers $headers -Body $topicsBody -ContentType "application/json" | Out-Null
        Write-Host "Topics added successfully!" -ForegroundColor Green
    } catch {
        Write-Host "Warning: Could not add topics automatically. Add them manually on GitHub." -ForegroundColor Yellow
    }
    
    # Add remote and push
    Write-Host ""
    Write-Host "Adding remote and pushing code..." -ForegroundColor Cyan
    git remote add origin "https://github.com/$Username/$RepoName.git" 2>$null
    if ($LASTEXITCODE -ne 0) {
        git remote set-url origin "https://github.com/$Username/$RepoName.git"
    }
    
    git push -u origin main
    if ($LASTEXITCODE -eq 0) {
        Write-Host ""
        Write-Host "Success! Repository is now on GitHub:" -ForegroundColor Green
        Write-Host $response.html_url -ForegroundColor Cyan
    }
} catch {
    Write-Host "Error creating repository: $_" -ForegroundColor Red
    Write-Host ""
    Write-Host "Please create it manually at https://github.com/new" -ForegroundColor Yellow
}

