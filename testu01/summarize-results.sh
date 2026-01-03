#!/usr/bin/env bash
# Summarize TestU01 results from a single log file
#
# Usage: ./summarize-results.sh <logfile>

if [ $# -eq 0 ]; then
    echo "Usage: $0 <logfile>"
    echo ""
    echo "Example:"
    echo "  $0 bigcrush-2.log"
    exit 1
fi

LOGFILE="$1"

if [ ! -f "$LOGFILE" ]; then
    echo "Error: File $LOGFILE not found"
    exit 1
fi

# Extract test suite name from filename (smallcrush, crush, or bigcrush)
TEST_TYPE=$(basename "$LOGFILE" | sed 's/-[0-9]*\.log$//' | sed 's/\.log$//')
TEST_NAME=$(echo "$TEST_TYPE" | tr '[:lower:]' '[:upper:]')

echo "═══════════════════════════════════════════════════════════"
echo "  TestU01 Results Summary"
echo "═══════════════════════════════════════════════════════════"
echo ""
echo "Log file: $LOGFILE"
echo "Test suite: $TEST_NAME"
echo "Generated: $(date)"
echo ""

# Extract test results
test_name=""
tests_found=0
passed=0
borderline=0
failed=0

# Create temporary file for results
tmpfile=$(mktemp)

# Parse the log file
subtest_name=""
while IFS= read -r line; do
    # Check for test name (lines ending with " test:")
    if echo "$line" | grep -q " test:$"; then
        test_name=$(echo "$line" | sed 's/ test://g' | xargs)
        subtest_name=""
    fi

    # Check for subtest (like "Test on the values of the Statistic H")
    if echo "$line" | grep -q "Test on the values of the Statistic"; then
        stat=$(echo "$line" | awk '{print $NF}')
        subtest_name=" ($stat)"
    fi

    # Check for p-value
    if echo "$line" | grep -q "p-value of test"; then
        pvalue=$(echo "$line" | awk '{print $NF}')

        if [ -n "$test_name" ] && [ -n "$pvalue" ]; then
            tests_found=$((tests_found + 1))

            # Determine status
            # Convert scientific notation to decimal (simple approach)
            pval_float="$pvalue"
            if echo "$pvalue" | grep -q "e-"; then
                # Simple conversion for common cases
                pval_float=$(echo "$pvalue" | sed '
                    s/1\.0e-3/0.001/g;
                    s/1\.0e-4/0.0001/g;
                    s/1\.0e-5/0.00001/g;
                ')
            fi

            # Determine status: FAIL < 0.001, BORDERLINE < 0.01, PASS otherwise
            if echo "$pval_float" | awk '{exit !($1 < 0.001 || $1 > 0.999)}'; then
                status="❌ FAIL"
                failed=$((failed + 1))
            elif echo "$pval_float" | awk '{exit !($1 < 0.01 || $1 > 0.99)}'; then
                status="⚠️  BORDERLINE"
                borderline=$((borderline + 1))
            else
                status="✅ PASS"
                passed=$((passed + 1))
            fi

            # Format test name (shorten if needed)
            short_name=$(echo "$test_name" | sed 's/^[a-z]*_//g')
            full_name="${short_name}${subtest_name}"
            display_name=$(echo "$full_name" | cut -c1-50)

            # Save to temp file
            printf "%-50s %12s  %s\n" "$display_name" "$pvalue" "$status" >> "$tmpfile"

            # Don't reset test_name - keep it until we see a new test
            # This allows multiple p-values per test to be captured
            subtest_name=""
        fi
    fi
done < "$LOGFILE"

# Display results table
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "  Detailed Results"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

if [ $tests_found -gt 0 ]; then
    printf "%-50s %12s  %s\n" "Test" "p-value" "Status"
    printf "%-50s %12s  %s\n" "─────────────────────────────────────────────────" "────────────" "─────────────"
    cat "$tmpfile"
    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "  Summary"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
    echo "  Tests run:      $tests_found"
    echo "  Passed:         $passed ✅"
    if [ $borderline -gt 0 ]; then
        echo "  Borderline:     $borderline ⚠️"
    fi
    if [ $failed -gt 0 ]; then
        echo "  Failed:         $failed ❌"
    fi

    # Calculate success rate
    if [ $tests_found -gt 0 ]; then
        success_rate=$(awk "BEGIN {printf \"%.1f\", ($passed / $tests_found) * 100}")
        echo "  Success rate:   ${success_rate}%"
    fi

    echo ""
    echo "  Final verdict:  $([ $failed -eq 0 ] && echo "✅ ALL TESTS PASSED" || echo "❌ SOME TESTS FAILED")"
else
    echo "  No test results found in log file"
fi

rm -f "$tmpfile"

echo ""
echo "═══════════════════════════════════════════════════════════"
echo "  Legend"
echo "═══════════════════════════════════════════════════════════"
echo "  ✅ PASS        p-value ∈ [0.01, 0.99]"
echo "  ⚠️  BORDERLINE  p-value ∈ [0.001, 0.01) ∪ (0.99, 0.999]"
echo "  ❌ FAIL        p-value < 0.001 or > 0.999"
echo ""
