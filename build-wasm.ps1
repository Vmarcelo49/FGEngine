# FGEngine WebAssembly Build Script
# This script compiles the project to WebAssembly and creates necessary files

param(
    [string]$Target = "main",  # "main" or "editor"
    [switch]$Serve = $false,   # Launch wasmserve after build
    [switch]$Clean = $false    # Clean existing files before build
)

Write-Host "=== FGEngine WebAssembly Build Script ===" -ForegroundColor Cyan
Write-Host "Target: $Target" -ForegroundColor Yellow

# Clean existing files if requested
if ($Clean) {
    Write-Host "Cleaning existing WebAssembly files..." -ForegroundColor Red
    Remove-Item *.wasm, *.html, wasm_exec.js -ErrorAction SilentlyContinue
    $cleanedFiles = @()
    if (Test-Path "*.wasm") { $cleanedFiles += "*.wasm" }
    if (Test-Path "*.html") { $cleanedFiles += "*.html" }
    if (Test-Path "wasm_exec.js") { $cleanedFiles += "wasm_exec.js" }
    
    if ($cleanedFiles.Count -gt 0) {
        Write-Host "Removed: $($cleanedFiles -join ', ')" -ForegroundColor Yellow
    } else {
        Write-Host "No files to clean" -ForegroundColor Gray
    }
}

# Set WebAssembly environment variables
Write-Host "Setting WebAssembly environment..." -ForegroundColor Green
$Env:GOOS = 'js'
$Env:GOARCH = 'wasm'

try {
    # Build the specified target
    if ($Target -eq "editor") {
        Write-Host "Building editor for WebAssembly..." -ForegroundColor Yellow
        go build -o editor.wasm ./cmd/editor
        if ($LASTEXITCODE -ne 0) {
            throw "Failed to build editor"
        }
        $wasmFile = "editor.wasm"
        $htmlFile = "editor.html"
        $title = "FGEngine Editor"
    }
    else {
        Write-Host "Building main application for WebAssembly..." -ForegroundColor Yellow
        go build -o main.wasm .
        if ($LASTEXITCODE -ne 0) {
            throw "Failed to build main application"
        }
        $wasmFile = "main.wasm"
        $htmlFile = "index.html"
        $title = "FGEngine"
    }

    # Get Go root for wasm_exec.js
    $goroot = go env GOROOT
    
    # Copy wasm_exec.js (try both locations for different Go versions)
    Write-Host "Copying wasm_exec.js..." -ForegroundColor Yellow
    $wasmExecPath = ""
    
    # Try Go 1.24+ location first
    if (Test-Path "$goroot\lib\wasm\wasm_exec.js") {
        $wasmExecPath = "$goroot\lib\wasm\wasm_exec.js"
    }
    # Fall back to Go 1.23 and older location
    elseif (Test-Path "$goroot\misc\wasm\wasm_exec.js") {
        $wasmExecPath = "$goroot\misc\wasm\wasm_exec.js"
    }
    else {
        Write-Host "Warning: wasm_exec.js not found in standard locations" -ForegroundColor Red
        Write-Host "You may need to download it manually from Go repository" -ForegroundColor Red
    }
    
    if ($wasmExecPath -ne "") {
        Copy-Item $wasmExecPath . -Force
        Write-Host "wasm_exec.js copied successfully" -ForegroundColor Green
    }

    # Create HTML file
    Write-Host "Creating HTML file: $htmlFile" -ForegroundColor Yellow
    
    $htmlContent = @"
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>$title</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        html, body {
            width: 100%;
            height: 100%;
            overflow: hidden;
            background-color: #000;
        }
        
        canvas {
            display: block;
            width: 100vw !important;
            height: 100vh !important;
            object-fit: contain;
        }
        
        .status {
            position: fixed;
            top: 50%;
            left: 50%;
            transform: translate(-50%, -50%);
            padding: 20px;
            border-radius: 8px;
            font-family: Arial, sans-serif;
            font-size: 16px;
            z-index: 1000;
            text-align: center;
            min-width: 200px;
        }
        
        .loading {
            background-color: rgba(0, 50, 100, 0.9);
            color: #ffffff;
            border: 2px solid #0080ff;
        }
        
        .error {
            background-color: rgba(100, 0, 0, 0.9);
            color: #ffffff;
            border: 2px solid #ff4444;
        }
        
        .ready {
            background-color: rgba(0, 100, 0, 0.9);
            color: #ffffff;
            border: 2px solid #44ff44;
        }
    </style>
</head>
<body>
    <div class="status loading" id="status">Loading WebAssembly...</div>

    <script src="wasm_exec.js"></script>
    <script>
        const statusElement = document.getElementById('status');
        
        function updateStatus(message, type = 'loading') {
            statusElement.textContent = message;
            statusElement.className = 'status ' + type;
        }
        
        async function loadWasm() {
            try {
                updateStatus('Initializing Go runtime...');
                const go = new Go();
                
                updateStatus('Fetching WebAssembly module...');
                const result = await WebAssembly.instantiateStreaming(fetch("$wasmFile"), go.importObject);
                
                updateStatus('Starting application...', 'ready');
                go.run(result.instance);
                
                // Hide status after successful load
                setTimeout(() => {
                    statusElement.style.display = 'none';
                }, 2000);
                
            } catch (error) {
                console.error('Failed to load WebAssembly:', error);
                updateStatus('Failed to load: ' + error.message, 'error');
            }
        }
        
        // Prevent context menu on right click
        document.addEventListener('contextmenu', function(e) {
            e.preventDefault();
        });
        
        // Prevent scrolling and zooming
        document.addEventListener('wheel', function(e) {
            e.preventDefault();
        }, { passive: false });
        
        document.addEventListener('touchstart', function(e) {
            if (e.touches.length > 1) {
                e.preventDefault();
            }
        }, { passive: false });
        
        // Start loading when page is ready
        if (document.readyState === 'loading') {
            document.addEventListener('DOMContentLoaded', loadWasm);
        } else {
            loadWasm();
        }
    </script>
</body>
</html>
"@

    $htmlContent | Out-File -FilePath $htmlFile -Encoding UTF8 -Force
    Write-Host "HTML file created successfully" -ForegroundColor Green

    # Show build results
    Write-Host "`n=== Build Results ===" -ForegroundColor Cyan
    if (Test-Path $wasmFile) {
        $wasmSize = (Get-Item $wasmFile).Length
        $wasmSizeMB = [math]::Round($wasmSize / 1MB, 2)
        Write-Host "✓ $wasmFile - $wasmSizeMB MB" -ForegroundColor Green
    }
    
    if (Test-Path "wasm_exec.js") {
        Write-Host "✓ wasm_exec.js" -ForegroundColor Green
    }
    
    if (Test-Path $htmlFile) {
        Write-Host "✓ $htmlFile" -ForegroundColor Green
    }

    Write-Host "`n=== Usage Instructions ===" -ForegroundColor Cyan
    Write-Host "1. Start a local HTTP server:" -ForegroundColor White
    Write-Host "   go run github.com/hajimehoshi/wasmserve@latest ." -ForegroundColor Yellow
    Write-Host "2. Open browser to: http://localhost:8080/$htmlFile" -ForegroundColor White
    Write-Host "`nOr use Python:" -ForegroundColor White
    Write-Host "   python -m http.server 8080" -ForegroundColor Yellow
    Write-Host "   # Then open: http://localhost:8080/$htmlFile" -ForegroundColor White

    # Optionally start wasmserve
    if ($Serve) {
        Write-Host "`nStarting wasmserve..." -ForegroundColor Green
        go run github.com/hajimehoshi/wasmserve@latest .
    }

} catch {
    Write-Host "Build failed: $_" -ForegroundColor Red
    exit 1
} finally {
    # Clean up environment variables
    Remove-Item Env:GOOS -ErrorAction SilentlyContinue
    Remove-Item Env:GOARCH -ErrorAction SilentlyContinue
}

Write-Host "`nBuild completed successfully!" -ForegroundColor Green