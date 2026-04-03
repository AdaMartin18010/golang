#!/bin/bash
#
# Go Knowledge Base Quality Check Script
#
# Description: Comprehensive quality checker for go-knowledge-base markdown files
# Version: 1.0.0
# Author: Knowledge Base Team
#
# Usage:
#   ./quality-check.sh [options] [path]
#
# Options:
#   -h, --help          Show help message
#   -j, --json          Output results as JSON
#   -m, --min-size N    Set minimum file size in KB (default: 15)
#   -s, --strict        Strict mode - fail on any quality issue
#   -q, --quiet         Quiet mode - only show errors
#   -d, --dimensions    Show statistics per dimension
#   -o, --output FILE   Save report to file
#   --fix               Attempt to auto-fix minor issues
#
# Exit Codes:
#   0 - All checks passed
#   1 - Quality issues found
#   2 - Configuration or file error
#
# Examples:
#   ./quality-check.sh                           # Check all files
#   ./quality-check.sh -j -o report.json         # JSON output to file
#   ./quality-check.sh -m 10 -s                  # Strict mode with 10KB min
#   ./quality-check.sh 01-Formal-Theory          # Check specific directory

set -o pipefail

#-------------------------------------------------------------------------------
# Configuration
#-------------------------------------------------------------------------------

SCRIPT_VERSION="1.0.0"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
KNOWLEDGE_BASE_DIR="$(dirname "$SCRIPT_DIR")"

# Default settings
MIN_SIZE_KB=15
MIN_SIZE_BYTES=$((MIN_SIZE_KB * 1024))
JSON_OUTPUT=false
STRICT_MODE=false
QUIET_MODE=false
SHOW_DIMENSIONS=false
OUTPUT_FILE=""
AUTO_FIX=false
TARGET_PATH=""

# Quality thresholds
readonly S_LEVEL_SIZE=15360      # 15KB
readonly A_LEVEL_SIZE=10240      # 10KB
readonly B_LEVEL_SIZE=5120       # 5KB
readonly C_LEVEL_SIZE=2048       # 2KB

readonly S_LEVEL_VISUALS=3
readonly A_LEVEL_VISUALS=2
readonly B_LEVEL_VISUALS=1

readonly S_LEVEL_REFS=5
readonly A_LEVEL_REFS=3
readonly B_LEVEL_REFS=1

# Statistics counters
declare -A STATS
declare -A DIMENSION_STATS
declare -a FILES_NEEDING_IMPROVEMENT
declare -a FILES_PASSED
declare -a FILES_FAILED

TOTAL_FILES=0
TOTAL_SIZE=0
PASSED_FILES=0
FAILED_FILES=0
WARNING_FILES=0

#-------------------------------------------------------------------------------
# Colors and Formatting
#-------------------------------------------------------------------------------

if [[ -t 1 ]]; then
    readonly RED='\033[0;31m'
    readonly GREEN='\033[0;32m'
    readonly YELLOW='\033[1;33m'
    readonly BLUE='\033[0;34m'
    readonly CYAN='\033[0;36m'
    readonly MAGENTA='\033[0;35m'
    readonly BOLD='\033[1m'
    readonly NC='\033[0m' # No Color
else
    readonly RED=''
    readonly GREEN=''
    readonly YELLOW=''
    readonly BLUE=''
    readonly CYAN=''
    readonly MAGENTA=''
    readonly BOLD=''
    readonly NC=''
fi

#-------------------------------------------------------------------------------
# Helper Functions
#-------------------------------------------------------------------------------

print_header() {
    if [[ "$QUIET_MODE" == false ]]; then
        echo -e "\n${BOLD}${BLUE}══════════════════════════════════════════════════════════════════════════════${NC}"
        echo -e "${BOLD}${BLUE}  $1${NC}"
        echo -e "${BOLD}${BLUE}══════════════════════════════════════════════════════════════════════════════${NC}\n"
    fi
}

print_section() {
    if [[ "$QUIET_MODE" == false ]]; then
        echo -e "\n${BOLD}${CYAN}▶ $1${NC}"
        echo -e "${CYAN}$(printf '─%.0s' {1..70})${NC}"
    fi
}

print_success() {
    if [[ "$QUIET_MODE" == false ]]; then
        echo -e "${GREEN}✓${NC} $1"
    fi
}

print_warning() {
    if [[ "$QUIET_MODE" == false ]]; then
        echo -e "${YELLOW}⚠${NC} $1"
    fi
}

print_error() {
    echo -e "${RED}✗${NC} $1" >&2
}

print_info() {
    if [[ "$QUIET_MODE" == false ]]; then
        echo -e "${BLUE}ℹ${NC} $1"
    fi
}

print_metric() {
    if [[ "$QUIET_MODE" == false ]]; then
        printf "${CYAN}%-30s${NC} %s\n" "$1" "$2"
    fi
}

#-------------------------------------------------------------------------------
# Argument Parsing
#-------------------------------------------------------------------------------

show_help() {
    cat << 'EOF'
Go Knowledge Base Quality Check Script

USAGE:
    ./quality-check.sh [OPTIONS] [PATH]

ARGUMENTS:
    PATH                Path to check (default: entire knowledge base)

OPTIONS:
    -h, --help          Show this help message and exit
    -j, --json          Output results in JSON format
    -m, --min-size N    Set minimum file size in KB (default: 15)
    -s, --strict        Strict mode - exit with error on any issue
    -q, --quiet         Quiet mode - only show errors
    -d, --dimensions    Show statistics grouped by dimension
    -o, --output FILE   Save report to specified file
    --fix               Attempt to auto-fix minor issues (not implemented)
    -v, --version       Show version information

EXAMPLES:
    # Check all files with default settings
    ./quality-check.sh

    # Check specific directory
    ./quality-check.sh 01-Formal-Theory

    # Generate JSON report
    ./quality-check.sh -j -o quality-report.json

    # Strict check with custom minimum size
    ./quality-check.sh -s -m 10

    # Show dimension statistics
    ./quality-check.sh -d

EXIT CODES:
    0   All checks passed
    1   Quality issues found
    2   Configuration or file error

For more information, see QUALITY-STANDARDS.md
EOF
}

parse_arguments() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -v|--version)
                echo "quality-check.sh version $SCRIPT_VERSION"
                exit 0
                ;;
            -j|--json)
                JSON_OUTPUT=true
                shift
                ;;
            -m|--min-size)
                if [[ -n "$2" && "$2" =~ ^[0-9]+$ ]]; then
                    MIN_SIZE_KB="$2"
                    MIN_SIZE_BYTES=$((MIN_SIZE_KB * 1024))
                    shift 2
                else
                    print_error "Option --min-size requires a numeric argument"
                    exit 2
                fi
                ;;
            -s|--strict)
                STRICT_MODE=true
                shift
                ;;
            -q|--quiet)
                QUIET_MODE=true
                shift
                ;;
            -d|--dimensions)
                SHOW_DIMENSIONS=true
                shift
                ;;
            -o|--output)
                if [[ -n "$2" ]]; then
                    OUTPUT_FILE="$2"
                    shift 2
                else
                    print_error "Option --output requires a filename"
                    exit 2
                fi
                ;;
            --fix)
                AUTO_FIX=true
                shift
                ;;
            -*)
                print_error "Unknown option: $1"
                echo "Use --help for usage information"
                exit 2
                ;;
            *)
                TARGET_PATH="$1"
                shift
                ;;
        esac
    done
}

#-------------------------------------------------------------------------------
# File Analysis Functions
#-------------------------------------------------------------------------------

get_file_size() {
    local file="$1"
    if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "cygwin" || "$OSTYPE" == "win32" ]]; then
        # Windows
        stat -c %s "$file" 2>/dev/null || wc -c < "$file"
    else
        # Unix-like
        stat -f%z "$file" 2>/dev/null || stat -c%s "$file" 2>/dev/null || wc -c < "$file"
    fi
}

get_dimension_from_path() {
    local file="$1"
    local rel_path="${file#$KNOWLEDGE_BASE_DIR/}"
    echo "$rel_path" | cut -d'/' -f1 | cut -d'\\' -f1
}

count_code_blocks() {
    local file="$1"
    grep -c '```go' "$file" 2>/dev/null || echo 0
}

count_mermaid_diagrams() {
    local file="$1"
    grep -c '```mermaid' "$file" 2>/dev/null || echo 0
}

count_ascii_diagrams() {
    local file="$1"
    grep -cE "^[[:space:]]*[┌├└│─═║╔╠╚]" "$file" 2>/dev/null || echo 0
}

count_cross_references() {
    local file="$1"
    grep -cE '\[.*\]\(.*\.md.*\)' "$file" 2>/dev/null || echo 0
}

count_headers() {
    local file="$1"
    grep -cE '^#{1,6} ' "$file" 2>/dev/null || echo 0
}

check_has_toc() {
    local file="$1"
    grep -qiE "^#{1,2} [Tt]able of [Cc]ontents|^\[Tt]able of [Cc]ontents\]" "$file" 2>/dev/null
    return $?
}

determine_quality_level() {
    local size="$1"
    local has_formal="$2"
    local visuals="$3"
    
    if [[ $size -ge $S_LEVEL_SIZE && "$has_formal" == "true" && $visuals -ge $S_LEVEL_VISUALS ]]; then
        echo "S"
    elif [[ $size -ge $A_LEVEL_SIZE && $visuals -ge $A_LEVEL_VISUALS ]]; then
        echo "A"
    elif [[ $size -ge $B_LEVEL_SIZE && $visuals -ge $B_LEVEL_VISUALS ]]; then
        echo "B"
    elif [[ $size -ge $C_LEVEL_SIZE ]]; then
        echo "C"
    else
        echo "F"
    fi
}

#-------------------------------------------------------------------------------
# File Processing
#-------------------------------------------------------------------------------

analyze_file() {
    local file="$1"
    local filename
    filename=$(basename "$file")
    
    local size
    size=$(get_file_size "$file")
    local size_kb=$((size / 1024))
    
    local dimension
    dimension=$(get_dimension_from_path "$file")
    
    # Count metrics
    local code_blocks
    code_blocks=$(count_code_blocks "$file")
    local mermaid_diagrams
    mermaid_diagrams=$(count_mermaid_diagrams "$file")
    local ascii_diagrams
    ascii_diagrams=$(count_ascii_diagrams "$file")
    local total_visuals=$((mermaid_diagrams + ascii_diagrams))
    local cross_refs
    cross_refs=$(count_cross_references "$file")
    local headers
    headers=$(count_headers "$file")
    
    # Check for formal content
    local has_formal=false
    if grep -qiE "(^#{1,3}.*[Ff]ormal [Dd]efinition|^[Dd]efinition|## [Ff]ormal|Theorem|Proof:)" "$file" 2>/dev/null; then
        has_formal=true
    fi
    
    # Check for TOC
    local has_toc=false
    if check_has_toc "$file"; then
        has_toc=true
    fi
    
    # Determine quality level
    local quality_level
    quality_level=$(determine_quality_level "$size" "$has_formal" "$total_visuals")
    
    # Build result object (JSON)
    local result=""
    result="{"
    result="${result}\"file\":\"$filename\","
    result="${result}\"path\":\"${file#$KNOWLEDGE_BASE_DIR/}\","
    result="${result}\"dimension\":\"$dimension\","
    result="${result}\"size_bytes\":$size,"
    result="${result}\"size_kb\":$size_kb,"
    result="${result}\"quality_level\":\"$quality_level\","
    result="${result}\"metrics\":{"
    result="${result}\"code_blocks\":$code_blocks,"
    result="${result}\"mermaid_diagrams\":$mermaid_diagrams,"
    result="${result}\"ascii_diagrams\":$ascii_diagrams,"
    result="${result}\"total_visuals\":$total_visuals,"
    result="${result}\"cross_references\":$cross_refs,"
    result="${result}\"headers\":$headers"
    result="${result}},"
    result="${result}\"features\":{"
    result="${result}\"has_formal_content\":$has_formal,"
    result="${result}\"has_toc\":$has_toc"
    result="${result}},"
    
    # Quality checks
    local issues=()
    local warnings=()
    
    if [[ $size -lt $MIN_SIZE_BYTES ]]; then
        issues+=("size-below-${MIN_SIZE_KB}kb")
    fi
    
    if [[ "$has_formal" == "false" && $size -ge $S_LEVEL_SIZE ]]; then
        warnings+=("s-level-without-formal-content")
    fi
    
    if [[ $total_visuals -lt 1 && $size -ge $B_LEVEL_SIZE ]]; then
        warnings+=("no-visualizations")
    fi
    
    if [[ $code_blocks -lt 1 && $size -ge $B_LEVEL_SIZE ]]; then
        warnings+=("no-code-examples")
    fi
    
    if [[ $cross_refs -lt 1 && $size -ge $B_LEVEL_SIZE ]]; then
        warnings+=("no-cross-references")
    fi
    
    # Build issues array
    local issues_json="["
    local first=true
    for issue in "${issues[@]}"; do
        if [[ "$first" == true ]]; then
            first=false
        else
            issues_json="${issues_json},"
        fi
        issues_json="${issues_json}\"$issue\""
    done
    issues_json="${issues_json}]"
    
    # Build warnings array
    local warnings_json="["
    first=true
    for warning in "${warnings[@]}"; do
        if [[ "$first" == true ]]; then
            first=false
        else
            warnings_json="${warnings_json},"
        fi
        warnings_json="${warnings_json}\"$warning\""
    done
    warnings_json="${warnings_json}]"
    
    # Determine pass/fail
    local passed=true
    if [[ ${#issues[@]} -gt 0 ]]; then
        passed=false
    fi
    
    result="${result}\"issues\":$issues_json,"
    result="${result}\"warnings\":$warnings_json,"
    result="${result}\"passed\":$passed"
    result="${result}}"
    
    echo "$result"
}

process_files() {
    local target="${1:-$KNOWLEDGE_BASE_DIR}"
    
    if [[ ! -d "$target" ]]; then
        target="$KNOWLEDGE_BASE_DIR/$target"
    fi
    
    if [[ ! -d "$target" ]]; then
        print_error "Target path not found: $1"
        exit 2
    fi
    
    # Initialize dimension stats
    DIMENSION_STATS["01-Formal-Theory"]=0
    DIMENSION_STATS["02-Language-Design"]=0
    DIMENSION_STATS["03-Engineering-CloudNative"]=0
    DIMENSION_STATS["04-Technology-Stack"]=0
    DIMENSION_STATS["05-Application-Domains"]=0
    
    # Find all markdown files
    local files=()
    while IFS= read -r -d '' file; do
        # Skip certain files
        local basename
        basename=$(basename "$file")
        if [[ "$basename" == README.md || "$basename" == CHANGELOG.md ]]; then
            continue
        fi
        
        # Skip files in root directory for dimension stats
        local rel_path="${file#$KNOWLEDGE_BASE_DIR/}"
        if [[ "$rel_path" != */* && "$rel_path" != *\\* ]]; then
            continue
        fi
        
        files+=("$file")
    done < <(find "$target" -name "*.md" -type f -print0 2>/dev/null)
    
    TOTAL_FILES=${#files[@]}
    
    if [[ $TOTAL_FILES -eq 0 ]]; then
        print_error "No markdown files found in target path"
        exit 2
    fi
    
    # Process each file
    local results=()
    
    for file in "${files[@]}"; do
        local result
        result=$(analyze_file "$file")
        results+=("$result")
        
        # Update statistics
        local size
        size=$(echo "$result" | grep -o '"size_bytes":[0-9]*' | cut -d: -f2)
        TOTAL_SIZE=$((TOTAL_SIZE + size))
        
        local dimension
        dimension=$(echo "$result" | grep -o '"dimension":"[^"]*"' | cut -d'"' -f4)
        if [[ -n "$dimension" ]]; then
            DIMENSION_STATS["$dimension"]=$((${DIMENSION_STATS[$dimension]:-0} + 1))
        fi
        
        local passed
        passed=$(echo "$result" | grep -o '"passed":[a-z]*' | cut -d: -f2)
        
        if [[ "$passed" == "true" ]]; then
            PASSED_FILES=$((PASSED_FILES + 1))
            FILES_PASSED+=("$file")
        else
            FAILED_FILES=$((FAILED_FILES + 1))
            FILES_FAILED+=("$file")
            local filepath
            filepath=$(echo "$result" | grep -o '"path":"[^"]*"' | cut -d'"' -f4)
            FILES_NEEDING_IMPROVEMENT+=("$filepath")
        fi
        
        # Show progress in non-quiet mode
        if [[ "$QUIET_MODE" == false && "$JSON_OUTPUT" == false ]]; then
            local current=$((PASSED_FILES + FAILED_FILES))
            if [[ $((current % 10)) -eq 0 || $current -eq $TOTAL_FILES ]]; then
                printf "\r${BLUE}Processing: ${NC}%d/%d files..." "$current" "$TOTAL_FILES"
            fi
        fi
    done
    
    if [[ "$QUIET_MODE" == false && "$JSON_OUTPUT" == false ]]; then
        echo -e "\n"
    fi
    
    # Output results
    if [[ "$JSON_OUTPUT" == true ]]; then
        output_json "${results[@]}"
    else
        output_text "${results[@]}"
    fi
}

#-------------------------------------------------------------------------------
# Output Formatting
#-------------------------------------------------------------------------------

output_json() {
    local results=("$@")
    
    local json="{"
    json="${json}\"meta\":{"
    json="${json}\"version\":\"$SCRIPT_VERSION\","
    json="${json}\"timestamp\":\"$(date -u +"%Y-%m-%dT%H:%M:%SZ")\","
    json="${json}\"target\":\"${TARGET_PATH:-$KNOWLEDGE_BASE_DIR}\","
    json="${json}\"min_size_kb\":$MIN_SIZE_KB"
    json="${json}},"
    
    json="${json}\"summary\":{"
    json="${json}\"total_files\":$TOTAL_FILES,"
    json="${json}\"passed\":$PASSED_FILES,"
    json="${json}\"failed\":$FAILED_FILES,"
    json="${json}\"pass_rate\":$(awk "BEGIN {printf \"%.2f\", ($PASSED_FILES/$TOTAL_FILES)*100}")"
    json="${json}},"
    
    # Dimension stats
    json="${json}\"dimensions\":{"
    local first=true
    for dim in "${!DIMENSION_STATS[@]}"; do
        if [[ "$first" == true ]]; then
            first=false
        else
            json="${json},"
        fi
        json="${json}\"$dim\":${DIMENSION_STATS[$dim]}"
    done
    json="${json}},"
    
    # Quality level distribution
    json="${json}\"quality_levels\":{"
    json="${json}\"S\":$(count_by_level "${results[@]}" "S"),"
    json="${json}\"A\":$(count_by_level "${results[@]}" "A"),"
    json="${json}\"B\":$(count_by_level "${results[@]}" "B"),"
    json="${json}\"C\":$(count_by_level "${results[@]}" "C"),"
    json="${json}\"F\":$(count_by_level "${results[@]}" "F")"
    json="${json}},"
    
    # Files needing improvement
    json="${json}\"needs_improvement\":["
    first=true
    for file in "${FILES_NEEDING_IMPROVMENT[@]}"; do
        if [[ "$first" == true ]]; then
            first=false
        else
            json="${json},"
        fi
        json="${json}\"$file\""
    done
    json="${json}],"
    
    # Detailed results
    json="${json}\"results\":["
    first=true
    for result in "${results[@]}"; do
        if [[ "$first" == true ]]; then
            first=false
        else
            json="${json},"
        fi
        json="${json}$result"
    done
    json="${json}]}"
    
    # Output
    if [[ -n "$OUTPUT_FILE" ]]; then
        echo "$json" | python3 -m json.tool 2>/dev/null || echo "$json" > "$OUTPUT_FILE"
        if [[ "$QUIET_MODE" == false ]]; then
            print_success "Report saved to: $OUTPUT_FILE"
        fi
    else
        echo "$json" | python3 -m json.tool 2>/dev/null || echo "$json"
    fi
}

count_by_level() {
    local results=("$@")
    local target_level="${results[-1]}"
    unset 'results[-1]'
    
    local count=0
    for result in "${results[@]}"; do
        local level
        level=$(echo "$result" | grep -oE '"quality_level":"[^"]*"' | cut -d'"' -f4)
        if [[ "$level" == "$target_level" ]]; then
            count=$((count + 1))
        fi
    done
    echo $count
}

output_text() {
    local results=("$@")
    
    print_header "QUALITY CHECK REPORT"
    
    # Meta info
    print_info "Version: $SCRIPT_VERSION"
    print_info "Timestamp: $(date -u +"%Y-%m-%dT%H:%M:%SZ")"
    print_info "Target: ${TARGET_PATH:-$KNOWLEDGE_BASE_DIR}"
    print_info "Minimum Size: ${MIN_SIZE_KB}KB"
    
    # Summary
    print_section "SUMMARY"
    print_metric "Total Files Checked:" "$TOTAL_FILES"
    print_metric "Files Passed:" "${GREEN}${PASSED_FILES}${NC}"
    print_metric "Files Failed:" "${RED}${FAILED_FILES}${NC}"
    
    local pass_rate
    pass_rate=$(awk "BEGIN {printf \"%.1f\", ($PASSED_FILES/$TOTAL_FILES)*100}")
    
    if [[ $FAILED_FILES -eq 0 ]]; then
        print_metric "Pass Rate:" "${GREEN}${pass_rate}%${NC} ${GREEN}✓ EXCELLENT${NC}"
    elif [[ $(awk "BEGIN {print ($pass_rate >= 80) ? 1 : 0}") -eq 1 ]]; then
        print_metric "Pass Rate:" "${YELLOW}${pass_rate}%${NC} ${YELLOW}⚠ GOOD${NC}"
    else
        print_metric "Pass Rate:" "${RED}${pass_rate}%${NC} ${RED}✗ NEEDS IMPROVEMENT${NC}"
    fi
    
    # Quality Level Distribution
    print_section "QUALITY LEVEL DISTRIBUTION"
    
    local s_count=0 a_count=0 b_count=0 c_count=0 f_count=0
    for result in "${results[@]}"; do
        local level
        level=$(echo "$result" | grep -oE '"quality_level":"[^"]*"' | cut -d'"' -f4)
        case "$level" in
            S) ((s_count++)) ;;
            A) ((a_count++)) ;;
            B) ((b_count++)) ;;
            C) ((c_count++)) ;;
            F) ((f_count++)) ;;
        esac
    done
    
    printf "  ${MAGENTA}S-Level${NC} (Supreme):     %3d files  ${BOLD}[>15KB + Formal + 3+ Visuals]${NC}\n" "$s_count"
    printf "  ${BLUE}A-Level${NC} (Advanced):    %3d files  ${BOLD}[>10KB + 2+ Visuals]${NC}\n" "$a_count"
    printf "  ${CYAN}B-Level${NC} (Basic):       %3d files  ${BOLD}[>5KB + 1+ Visual]${NC}\n" "$b_count"
    printf "  ${YELLOW}C-Level${NC} (Concise):     %3d files  ${BOLD}[>2KB]${NC}\n" "$c_count"
    printf "  ${RED}F-Level${NC} (Failed):      %3d files  ${BOLD}[<2KB or missing requirements]${NC}\n" "$f_count"
    
    # Dimension Statistics
    if [[ "$SHOW_DIMENSIONS" == true ]]; then
        print_section "STATISTICS BY DIMENSION"
        
        for dim in "01-Formal-Theory" "02-Language-Design" "03-Engineering-CloudNative" "04-Technology-Stack" "05-Application-Domains"; do
            local count=${DIMENSION_STATS[$dim]:-0}
            local dim_name
            dim_name=$(echo "$dim" | sed 's/^[0-9]*-//' | tr '-' ' ')
            printf "  ${CYAN}%-30s${NC} %3d files\n" "$dim_name:" "$count"
        done
    fi
    
    # Files Needing Improvement
    if [[ ${#FILES_NEEDING_IMPROVEMENT[@]} -gt 0 ]]; then
        print_section "FILES NEEDING IMPROVEMENT"
        
        for file in "${FILES_NEEDING_IMPROVEMENT[@]}"; do
            print_error "$file"
        done
        
        if [[ ${#FILES_NEEDING_IMPROVEMENT[@]} -gt 10 ]]; then
            local remaining
            remaining=$((${#FILES_NEEDING_IMPROVEMENT[@]} - 10))
            print_info "... and $remaining more files (use -j for full list)"
        fi
    fi
    
    # Recommendations
    print_section "RECOMMENDATIONS"
    
    if [[ $FAILED_FILES -eq 0 ]]; then
        print_success "All files meet the minimum quality standards!"
    else
        if [[ $f_count -gt 0 ]]; then
            print_error "$f_count files are below minimum quality threshold (2KB)"
            print_info "  → Consider expanding content or merging with related documents"
        fi
        
        if [[ ${#FILES_NEEDING_IMPROVEMENT[@]} -gt 0 ]]; then
            print_warning "${#FILES_NEEDING_IMPROVEMENT[@]} files need improvements"
            print_info "  → Add missing sections: definitions, theorems, code examples, visualizations"
            print_info "  → See QUALITY-STANDARDS.md for detailed requirements"
        fi
    fi
    
    # Target S-Level percentage
    local s_percentage
    s_percentage=$(awk "BEGIN {printf \"%.1f\", ($s_count/$TOTAL_FILES)*100}")
    print_info "Current S-Level percentage: ${s_percentage}% (target: 15%)"
    
    # Footer
    echo -e "\n${BOLD}${BLUE}══════════════════════════════════════════════════════════════════════════════${NC}"
    echo -e "${BLUE}For more details, see:${NC}"
    echo -e "  ${CYAN}• QUALITY-STANDARDS.md${NC} - Quality level definitions"
    echo -e "  ${CYAN}• TEMPLATES.md${NC} - Document templates and examples"
    echo -e "${BOLD}${BLUE}══════════════════════════════════════════════════════════════════════════════${NC}\n"
    
    # Save to file if requested
    if [[ -n "$OUTPUT_FILE" && "$JSON_OUTPUT" == false ]]; then
        {
            echo "QUALITY CHECK REPORT"
            echo "===================="
            echo "Version: $SCRIPT_VERSION"
            echo "Timestamp: $(date -u +"%Y-%m-%dT%H:%M:%SZ")"
            echo ""
            echo "SUMMARY"
            echo "-------"
            echo "Total Files: $TOTAL_FILES"
            echo "Passed: $PASSED_FILES"
            echo "Failed: $FAILED_FILES"
            echo "Pass Rate: ${pass_rate}%"
            echo ""
            echo "QUALITY LEVELS"
            echo "--------------"
            echo "S-Level: $s_count"
            echo "A-Level: $a_count"
            echo "B-Level: $b_count"
            echo "C-Level: $c_count"
            echo "F-Level: $f_count"
            echo ""
            echo "FILES NEEDING IMPROVEMENT"
            echo "-------------------------"
            printf '%s\n' "${FILES_NEEDING_IMPROVEMENT[@]}"
        } > "$OUTPUT_FILE"
        print_success "Report saved to: $OUTPUT_FILE"
    fi
}

#-------------------------------------------------------------------------------
# Main
#-------------------------------------------------------------------------------

main() {
    parse_arguments "$@"
    
    # Validate knowledge base directory
    if [[ ! -d "$KNOWLEDGE_BASE_DIR" ]]; then
        print_error "Knowledge base directory not found: $KNOWLEDGE_BASE_DIR"
        exit 2
    fi
    
    # Check dependencies
    if ! command -v grep &>/dev/null; then
        print_error "Required command 'grep' not found"
        exit 2
    fi
    
    # Run quality check
    if [[ "$JSON_OUTPUT" == false && "$QUIET_MODE" == false ]]; then
        print_header "Go Knowledge Base Quality Check"
        print_info "Scanning for markdown files..."
    fi
    
    process_files "$TARGET_PATH"
    
    # Exit with appropriate code
    if [[ $FAILED_FILES -gt 0 && "$STRICT_MODE" == true ]]; then
        exit 1
    elif [[ $FAILED_FILES -gt 0 ]]; then
        exit 1
    else
        exit 0
    fi
}

# Run main function
main "$@"
