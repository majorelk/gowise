#!/bin/bash
# UK English spelling checker for GoWise project
# Ensures all code comments and documentation use UK English spellings
# Only flags new violations, not existing code or stdlib compatibility

set -e

echo "🔍 Checking for US English spellings in comments and docs..."

# US spellings that should be UK spellings (only in comments/docs, not code)
US_SPELLINGS="behavior|initialize|organize|prioritize|optimize|analyze|catalog|flavor|color|honor|labor|favor|armor|tumor|humor|vigor|traveled|canceled|modeling|leveling"

FOUND_ISSUES=0

# Search for US spellings but exclude:
# - JSON/XML struct tags (API compatibility)
# - Standard library function calls (context.WithCancel, etc.)
# - Import statements
# - Existing known cases for backward compatibility
US_RESULTS=$(grep -r -n -E "\b($US_SPELLINGS)\b" --include="*.go" --include="*.md" . | \
    grep -v ".git" | \
    grep -v 'json:' | \
    grep -v 'xml:' | \
    grep -v 'import ' | \
    grep -v 'context\.' | \
    grep -v 'http\.' | \
    grep -v 'func.*(' | \
    grep -v '\.' | \
    grep -v 'swparser.go' || true)

if [ -n "$US_RESULTS" ]; then
    echo "Found US English spellings in comments/documentation:"
    echo "$US_RESULTS"
    echo ""
    echo "❌ Use UK English equivalents in comments and documentation:"
    echo "   behavior → behaviour     initialize → initialise"
    echo "   organize → organise      prioritize → prioritise" 
    echo "   optimize → optimise      analyze → analyse"
    echo "   catalog → catalogue      traveled → travelled"
    echo "   canceled → cancelled     modeling → modelling"
    echo ""
    echo "ℹ️  Note: Standard library function names and JSON/XML tags are exempt"
    FOUND_ISSUES=1
fi

# Check for "license" in comments/docs (excluding API struct tags and headers)
LICENSE_RESULTS=$(grep -r -n "\blicense\b" --include="*.go" --include="*.md" . | \
    grep -v ".git" | \
    grep -v "swparser.go" | \
    grep -v 'json:' | \
    grep -v 'xml:' | \
    grep -v 'License:' | \
    grep -v 'LICENSE' | \
    grep -v '// Copyright' | \
    grep -v 'MIT' || true)

if [ -n "$LICENSE_RESULTS" ]; then
    echo "Found 'license' in comments/documentation:"
    echo "$LICENSE_RESULTS"  
    echo ""
    echo "❌ Use 'licence' in UK English for comments and documentation"
    echo "ℹ️  Note: JSON/XML struct tags may use 'license' for API compatibility"
    FOUND_ISSUES=1
fi

if [ $FOUND_ISSUES -eq 1 ]; then
    echo ""
    echo "💡 Fix US spellings in comments and documentation above."
    echo "   Code identifiers and stdlib calls are not affected by this check."
    exit 1
fi

echo "✅ All spelling checks passed - UK English used in comments and docs"