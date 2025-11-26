# PowerShell script to create GitHub repository and push code
# Make sure you're logged into GitHub and have a personal access token

$repoName = "http-cli"
$description = "A colorful, lightweight command-line HTTP client similar to Postman. Built with Go for fast performance and beautiful terminal output."
$tags = "go,cli,http-client,rest-api,postman-alternative,terminal-tool,http-requests,api-testing,command-line,developer-tools"

Write-Host "Creating GitHub repository: $repoName" -ForegroundColor Cyan
Write-Host "Description: $description" -ForegroundColor Gray
Write-Host ""
Write-Host "Please follow these steps:" -ForegroundColor Yellow
Write-Host "1. Go to https://github.com/new" -ForegroundColor White
Write-Host "2. Repository name: $repoName" -ForegroundColor White
Write-Host "3. Description: $description" -ForegroundColor White
Write-Host "4. Set to Public (or Private if you prefer)" -ForegroundColor White
Write-Host "5. DO NOT initialize with README, .gitignore, or license" -ForegroundColor White
Write-Host "6. Click 'Create repository'" -ForegroundColor White
Write-Host ""
Write-Host "After creating the repository, run these commands:" -ForegroundColor Yellow
Write-Host ""
Write-Host "git remote add origin https://github.com/YOUR_USERNAME/$repoName.git" -ForegroundColor Green
Write-Host "git branch -M main" -ForegroundColor Green
Write-Host "git push -u origin main" -ForegroundColor Green
Write-Host ""
Write-Host "Then add topics/tags on GitHub: $tags" -ForegroundColor Cyan

