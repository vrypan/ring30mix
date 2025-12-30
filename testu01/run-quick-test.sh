#!/bin/bash
# Quick test runner - checks if TestU01 is set up and runs SmallCrush

set -e

echo "═══════════════════════════════════════════════════════════"
echo "  Rule30 RNG - Quick TestU01 Test"
echo "═══════════════════════════════════════════════════════════"
echo ""

# Check if TestU01 is installed
if [ ! -f /usr/local/lib/libtestu01.a ]; then
    echo "TestU01 is not installed."
    echo ""
    read -p "Install TestU01 now? (y/N) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        ./install-testu01.sh
    else
        echo "Installation cancelled. Please install TestU01 manually."
        echo "See README.md for instructions."
        exit 1
    fi
fi

# Check if test programs are built
if [ ! -f ./test-smallcrush ]; then
    echo "Building test programs..."
    make all
    echo "✓ Test programs built"
    echo ""
fi

# Check if rule30 binary exists
if [ ! -f ../rule30 ]; then
    echo "Building rule30 binary..."
    cd .. && make rule30 && cd testu01
    echo "✓ rule30 binary built"
    echo ""
fi

# Run SmallCrush
echo "Running SmallCrush test..."
echo "This will take approximately 1-2 minutes."
echo ""
echo "Random data will be piped from Rule30 RNG to TestU01."
echo "Look for 'All tests were passed' or individual test failures."
echo ""
echo "Press Enter to start..."
read

cd .. && ./rule30 --bytes=524288000 2>&1 | ./testu01/test-smallcrush 2>&1 | grep -v "Error: Failed"

echo ""
echo "═══════════════════════════════════════════════════════════"
echo "  Test Complete"
echo "═══════════════════════════════════════════════════════════"
echo ""
echo "Next steps:"
echo "  - Run full Crush:    cd testu01 && make crush"
echo "  - Run BigCrush:      cd testu01 && make bigcrush"
echo ""
