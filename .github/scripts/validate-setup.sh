#!/bin/bash

# GitHub Actions Setup Validation Script
# This script helps validate your GitHub Actions setup

echo "🔍 Validating GitHub Actions setup..."

# Check if workflow files exist
WORKFLOWS_DIR=".github/workflows"
if [ ! -d "$WORKFLOWS_DIR" ]; then
    echo "❌ GitHub workflows directory not found"
    exit 1
fi

echo "✅ GitHub workflows directory exists"

# Check workflow files
REQUIRED_WORKFLOWS=("ci.yml" "deploy.yml" "release.yml")
for workflow in "${REQUIRED_WORKFLOWS[@]}"; do
    if [ -f "$WORKFLOWS_DIR/$workflow" ]; then
        echo "✅ $workflow found"
    else
        echo "❌ $workflow missing"
        exit 1
    fi
done

# Check Go files
if [ ! -f "go.mod" ]; then
    echo "❌ go.mod not found"
    exit 1
fi
echo "✅ go.mod found"

if [ ! -f "main.go" ]; then
    echo "❌ main.go not found"
    exit 1
fi
echo "✅ main.go found"

# Test Go build
echo "🔨 Testing Go build..."
if go build -v .; then
    echo "✅ Go build successful"
else
    echo "❌ Go build failed"
    exit 1
fi

# Check for required environment variables (optional)
echo "🔑 Checking environment variables..."
REQUIRED_VARS=("TRADING212_API_KEY" "CMC_API_KEY" "BULLION_API_KEY")
for var in "${REQUIRED_VARS[@]}"; do
    if [ -z "${!var}" ]; then
        echo "⚠️  $var not set (will need to be set in GitHub secrets)"
    else
        echo "✅ $var is set"
    fi
done

echo ""
echo "🎉 GitHub Actions setup validation complete!"
echo ""
echo "Next steps:"
echo "1. Push your code to GitHub"
echo "2. Set up repository secrets (see .github/SETUP.md)"
echo "3. Watch the CI workflow run automatically"
echo "4. Manually trigger the deploy workflow to test"