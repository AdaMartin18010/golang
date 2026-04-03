#!/bin/bash

# Go Knowledge Base - Automatic Index Generator
# This script scans all markdown files and generates cross-referenced indices
# Usage: ./generate-index.sh [knowledge_base_path]
# Default path: ../ (parent of scripts directory)

set -euo pipefail

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
KB_PATH="${1:-$SCRIPT_DIR/..}"
KB_PATH="$(cd "$KB_PATH" && pwd)"
INDICES_DIR="$KB_PATH/indices"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() { echo -e "${BLUE}[INFO]${NC} $1"; }
log_success() { echo -e "${GREEN}[SUCCESS]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# Statistics
STATS_TOTAL=0
STATS_WITH_METADATA=0
STATS_DIMENSIONS=""
STATS_LEVELS=""

# Create indices directory if it doesn't exist
ensure_directories() {
    mkdir -p "$INDICES_DIR"
    log_info "Indices directory: $INDICES_DIR"
}

# Extract metadata from a markdown file
# Returns: title|dimension|category|level|tags|date|path
extract_metadata() {
    local file="$1"
    local rel_path="${file#$KB_PATH/}"
    local filename=$(basename "$file")
    
    # Read first 50 lines to extract metadata
    local header=$(head -50 "$file" 2>/dev/null || echo "")
    
    # Extract title (first # line)
    local title=$(echo "$header" | grep -m1 '^# ' | sed 's/^# //' | head -c 100)
    if [[ -z "$title" ]]; then
        title="$filename"
    fi
    
    # Extract dimension (维度)
    local dimension=$(echo "$header" | grep -oP '\*\*维度\*\*[:：]\s*\K[^|]+' | head -1 | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
    if [[ -z "$dimension" ]]; then
        # Try alternative patterns
        dimension=$(echo "$header" | grep -oP '>\s*\*\*维度\*\*\s*[:|]\s*\K[^|)]+' | head -1 | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
    fi
    
    # Extract category (分类)
    local category=$(echo "$header" | grep -oP '\*\*分类\*\*[:：]\s*\K[^|]+' | head -1 | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
    
    # Extract level (级别/难度)
    local level=""
    if echo "$header" | grep -qP '\*\*级别\*\*[:：]\s*\K[SABC]'; then
        level=$(echo "$header" | grep -oP '\*\*级别\*\*[:：]\s*\K[SABC][^[:space:]]*' | head -1)
    elif echo "$header" | grep -qP '\*\*难度\*\*[:：]\s*\K(初级|中级|高级)'; then
        level=$(echo "$header" | grep -oP '\*\*难度\*\*[:：]\s*\K[^[:space:]]+' | head -1)
        # Normalize level
        case "$level" in
            初级|Beginner) level="B" ;;
            中级|Intermediate) level="A" ;;
            高级|Advanced|Expert) level="S" ;;
        esac
    elif echo "$header" | grep -iq 'S-Level\|S级'; then
        level="S"
    elif echo "$header" | grep -iq 'A-Level\|A级'; then
        level="A"
    elif echo "$header" | grep -iq 'B-Level\|B级'; then
        level="B"
    fi
    
    # Extract tags
    local tags=""
    if echo "$header" | grep -qP '\*\*标签\*\*[:：]'; then
        tags=$(echo "$header" | grep -oP '\*\*标签\*\*[:：]\s*\K[^|]+' | head -1 | tr ',' '|' | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
    elif echo "$header" | grep -qP '#[a-z-]+'; then
        tags=$(echo "$header" | grep -oP '#\K[a-z-]+' | tr '\n' '|' | sed 's/|$//')
    fi
    
    # Extract date
    local date=""
    if echo "$header" | grep -qP '\*\*最后更新\*\*[:：]'; then
        date=$(echo "$header" | grep -oP '\*\*最后更新\*\*[:：]\s*\K[0-9]{4}-[0-9]{2}-[0-9]{2}' | head -1)
    elif echo "$header" | grep -qP '\*\*完成日期\*\*[:：]'; then
        date=$(echo "$header" | grep -oP '\*\*完成日期\*\*[:：]\s*\K[0-9]{4}-[0-9]{2}-[0-9]{2}' | head -1)
    elif echo "$header" | grep -qP '\*\*Created\*\*[:：]'; then
        date=$(echo "$header" | grep -oP '\*\*Created\*\*[:：]\s*\K[0-9]{4}-[0-9]{2}-[0-9]{2}' | head -1)
    fi
    
    # Extract file modification date as fallback
    if [[ -z "$date" ]]; then
        date=$(stat -c %Y "$file" 2>/dev/null | xargs -I{} date -d @{} +%Y-%m-%d 2>/dev/null || stat -f %Sm -t %Y-%m-%d "$file" 2>/dev/null || echo "")
    fi
    
    # Output delimited string
    echo "$title|$dimension|$category|$level|$tags|$date|$rel_path"
}

# Get dimension from path
get_dimension_from_path() {
    local path="$1"
    local dim=""
    
    if [[ "$path" =~ ^01-Formal-Theory ]]; then
        dim="Formal Theory"
    elif [[ "$path" =~ ^02-Language-Design ]]; then
        dim="Language Design"
    elif [[ "$path" =~ ^03-Engineering-CloudNative ]]; then
        dim="Engineering & Cloud Native"
    elif [[ "$path" =~ ^04-Technology-Stack ]]; then
        dim="Technology Stack"
    elif [[ "$path" =~ ^05-Application-Domains ]]; then
        dim="Application Domains"
    elif [[ "$path" =~ ^examples ]]; then
        dim="Examples"
    elif [[ "$path" =~ ^indices ]]; then
        dim="Indices"
    elif [[ "$path" =~ ^learning-paths ]]; then
        dim="Learning Paths"
    else
        dim="Other"
    fi
    
    echo "$dim"
}

# Get level from file size and content
get_level_heuristic() {
    local file="$1"
    local size=$(stat -f%z "$file" 2>/dev/null || stat -c%s "$file" 2>/dev/null || echo "0")
    
    if [[ $size -gt 15000 ]]; then
        echo "S"
    elif [[ $size -gt 10000 ]]; then
        echo "A"
    elif [[ $size -gt 5000 ]]; then
        echo "B"
    else
        echo "C"
    fi
}

# Scan all markdown files
scan_files() {
    log_info "Scanning markdown files in $KB_PATH..."
    
    local temp_file=$(mktemp)
    local count=0
    
    # Find all markdown files excluding indices and certain directories
    while IFS= read -r -d '' file; do
        # Skip certain files
        local rel_path="${file#$KB_PATH/}"
        if [[ "$rel_path" =~ ^(indices|scripts)/ ]]; then
            continue
        fi
        
        local metadata=$(extract_metadata "$file")
        
        # Parse metadata
        local title=$(echo "$metadata" | cut -d'|' -f1)
        local dimension=$(echo "$metadata" | cut -d'|' -f2)
        local category=$(echo "$metadata" | cut -d'|' -f3)
        local level=$(echo "$metadata" | cut -d'|' -f4)
        local tags=$(echo "$metadata" | cut -d'|' -f5)
        local date=$(echo "$metadata" | cut -d'|' -f6)
        local path=$(echo "$metadata" | cut -d'|' -f7)
        
        # Infer dimension from path if not found
        if [[ -z "$dimension" ]]; then
            dimension=$(get_dimension_from_path "$path")
        fi
        
        # Infer level from file size if not found
        if [[ -z "$level" ]]; then
            level=$(get_level_heuristic "$file")
        fi
        
        # Use filename as date fallback
        if [[ -z "$date" ]]; then
            date="$(date +%Y-%m-%d)"
        fi
        
        # Write to temp file
        echo "$title|$dimension|$category|$level|$tags|$date|$path" >> "$temp_file"
        
        ((count++))
        if (( count % 100 == 0 )); then
            echo -ne "\r${BLUE}[INFO]${NC} Scanned $count files..."
        fi
        
    done < <(find "$KB_PATH" -type f -name "*.md" -print0 2>/dev/null | grep -zEv '/(indices|scripts)/')
    
    echo ""
    STATS_TOTAL=$count
    log_success "Scanned $count markdown files"
    
    echo "$temp_file"
}

# Generate by-topic.md
generate_by_topic() {
    local temp_file="$1"
    local output="$INDICES_DIR/by-topic.md"
    
    log_info "Generating by-topic.md..."
    
    cat > "$output" << 'EOF'
# Go Knowledge Base - Topic-Based Cross Reference

> **Version**: Auto-generated
> **Last Updated**: 
EOF
    echo "$(date +%Y-%m-%d)" >> "$output"
    cat >> "$output" << 'EOF'
> **Purpose**: Find documents by topic, concept, or technology

---

## 📑 Topics Overview

EOF

    # Build topic index from dimensions and categories
    declare -A topics
    
    while IFS='|' read -r title dimension category level tags date path; do
        # Add dimension as topic
        if [[ -n "$dimension" ]]; then
            topics["$dimension"]+="$path|$title|dimension\n"
        fi
        
        # Add category as topic
        if [[ -n "$category" ]]; then
            topics["$category"]+="$path|$title|category\n"
        fi
        
        # Add tags as topics
        if [[ -n "$tags" ]]; then
            IFS='|' read -ra tag_array <<< "$tags"
            for tag in "${tag_array[@]}"; do
                tag=$(echo "$tag" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
                if [[ -n "$tag" ]]; then
                    topics["$tag"]+="$path|$title|tag\n"
                fi
            done
        fi
        
        # Extract topics from path
        if [[ "$path" =~ ^[0-9]+-([^/]+) ]]; then
            local dim_topic="${BASH_REMATCH[1]}"
            dim_topic=$(echo "$dim_topic" | sed 's/-/ /g' | sed 's/\b\w/\u&/g')
            topics["$dim_topic"]+="$path|$title|path\n"
        fi
    done < "$temp_file"
    
    # Sort topics and generate alphabetically
    local sorted_topics=$(printf '%s\n' "${!topics[@]}" | sort)
    
    # Group by first letter
    local current_letter=""
    for topic in $sorted_topics; do
        local first_letter=$(echo "$topic" | head -c 1 | tr '[:lower:]' '[:upper:]')
        if [[ "$first_letter" != "$current_letter" ]]; then
            current_letter="$first_letter"
            echo "" >> "$output"
            echo "### $current_letter" >> "$output"
            echo "" >> "$output"
        fi
        
        echo "- **$topic**" >> "$output"
        
        # List documents under this topic (limit to 5)
        local docs="${topics[$topic]}"
        local count=0
        while IFS='|' read -r doc_path doc_title doc_type; do
            if [[ -n "$doc_path" && $count -lt 5 ]]; then
                local short_path=$(echo "$doc_path" | sed 's/\.md$//')
                echo "  - [$doc_title](../$doc_path)" >> "$output"
                ((count++))
            fi
        done <<< "$docs"
        
        if [[ $(echo "$docs" | grep -c "|") -gt 5 ]]; then
            echo "  - *(and more...)*" >> "$output"
        fi
        echo "" >> "$output"
    done
    
    cat >> "$output" << 'EOF'

---

*This index is automatically generated. Run `./scripts/generate-index.sh` to update.*
EOF

    log_success "Generated $output ($(wc -l < "$output") lines)"
}

# Generate by-tag.md
generate_by_tag() {
    local temp_file="$1"
    local output="$INDICES_DIR/by-tag.md"
    
    log_info "Generating by-tag.md..."
    
    cat > "$output" << 'EOF'
# Go Knowledge Base - Tag Index

> **Version**: Auto-generated
> **Last Updated**: 
EOF
    echo "$(date +%Y-%m-%d)" >> "$output"
    cat >> "$output" << 'EOF'
> **Purpose**: Find documents by tags

---

## 🏷️ Tag Index

EOF

    # Collect tags
    declare -A tag_docs
    
    while IFS='|' read -r title dimension category level tags date path; do
        if [[ -n "$tags" ]]; then
            IFS='|' read -ra tag_array <<< "$tags"
            for tag in "${tag_array[@]}"; do
                tag=$(echo "$tag" | sed 's/^[#[:space:]]*//;s/[[:space:]]*$//' | tr '[:upper:]' '[:lower:]')
                if [[ -n "$tag" ]]; then
                    tag_docs["$tag"]+="- [$title](../$path)\n"
                fi
            done
        fi
    done < "$temp_file"
    
    # Also use dimension and category as tags
    while IFS='|' read -r title dimension category level tags date path; do
        if [[ -n "$dimension" ]]; then
            local dim_tag=$(echo "$dimension" | tr '[:upper:]' '[:lower:]' | sed 's/[^a-z0-9]/-/g' | sed 's/--*/-/g' | sed 's/^-//;s/-$//')
            if [[ -n "$dim_tag" ]]; then
                tag_docs["$dim_tag"]+="- [$title](../$path)\n"
            fi
        fi
        if [[ -n "$category" ]]; then
            local cat_tag=$(echo "$category" | tr '[:upper:]' '[:lower:]' | sed 's/[^a-z0-9]/-/g' | sed 's/--*/-/g' | sed 's/^-//;s/-$//')
            if [[ -n "$cat_tag" && "$cat_tag" != "$dim_tag" ]]; then
                tag_docs["$cat_tag"]+="- [$title](../$path)\n"
            fi
        fi
    done < "$temp_file"
    
    # Sort tags and output
    local sorted_tags=$(printf '%s\n' "${!tag_docs[@]}" | sort)
    
    for tag in $sorted_tags; do
        local docs="${tag_docs[$tag]}"
        local count=$(echo -e "$docs" | grep -c "^-" || echo "0")
        
        echo "### #$tag ($count documents)" >> "$output"
        echo "" >> "$output"
        echo -e "$docs" | sort -u >> "$output"
        echo "" >> "$output"
    done
    
    cat >> "$output" << 'EOF'

---

*This index is automatically generated. Run `./scripts/generate-index.sh` to update.*
EOF

    log_success "Generated $output ($(wc -l < "$output") lines)"
}

# Generate by-date.md
generate_by_date() {
    local temp_file="$1"
    local output="$INDICES_DIR/by-date.md"
    
    log_info "Generating by-date.md..."
    
    cat > "$output" << 'EOF'
# Go Knowledge Base - Chronological Index

> **Version**: Auto-generated
> **Last Updated**: 
EOF
    echo "$(date +%Y-%m-%d)" >> "$output"
    cat >> "$output" << 'EOF'
> **Purpose**: Find documents by creation/update date

---

## 📅 Documents by Date

EOF

    # Collect by date
    declare -A date_docs
    
    while IFS='|' read -r title dimension category level tags date path; do
        if [[ -n "$date" ]]; then
            # Normalize date format
            local norm_date=$(echo "$date" | grep -oP '^[0-9]{4}-[0-9]{2}-[0-9]{2}' || echo "$date")
            date_docs["$norm_date"]+="|[$title](../$path)|$dimension|$level\n"
        fi
    done < "$temp_file"
    
    # Sort dates in reverse chronological order
    local sorted_dates=$(printf '%s\n' "${!date_docs[@]}" | sort -r)
    
    # Group by month
    local current_month=""
    for date in $sorted_dates; do
        local month=$(echo "$date" | cut -d'-' -f1,2)
        if [[ "$month" != "$current_month" ]]; then
            current_month="$month"
            echo "" >> "$output"
            echo "### $current_month" >> "$output"
            echo "" >> "$output"
            echo "| Date | Document | Dimension | Level |" >> "$output"
            echo "|------|----------|-----------|-------|" >> "$output"
        fi
        
        local docs="${date_docs[$date]}"
        while IFS='|' read -r doc dim lvl; do
            if [[ -n "$doc" ]]; then
                echo "| $date | $doc | $dim | $lvl |" >> "$output"
            fi
        done <<< "$docs"
    done
    
    echo "" >> "$output"
    
    # Add monthly statistics
    cat >> "$output" << 'EOF'

## 📊 Monthly Statistics

EOF

    for month in $(printf '%s\n' "${!date_docs[@]}" | cut -d'-' -f1,2 | sort -u -r); do
        local count=$(printf '%s\n' "${!date_docs[@]}" | grep "^$month" | wc -l)
        echo "- **$month**: $count documents" >> "$output"
    done
    
    cat >> "$output" << 'EOF'

---

*This index is automatically generated. Run `./scripts/generate-index.sh` to update.*
EOF

    log_success "Generated $output ($(wc -l < "$output") lines)"
}

# Generate complete-map.md
generate_complete_map() {
    local temp_file="$1"
    local output="$INDICES_DIR/complete-map.md"
    
    log_info "Generating complete-map.md..."
    
    local total_docs=$(wc -l < "$temp_file")
    local dim_stats=""
    local level_stats=""
    
    # Calculate dimension statistics
    dim_stats=$(cut -d'|' -f2 "$temp_file" | sort | uniq -c | sort -rn)
    level_stats=$(cut -d'|' -f4 "$temp_file" | sort | uniq -c | sort -rn)
    
    cat > "$output" << EOF
# Go Knowledge Base - Complete Document Map

> **Version**: Auto-generated
> **Last Updated**: $(date +%Y-%m-%d)
> **Purpose**: Complete inventory of all knowledge base documents

---

## 📊 Statistics

| Metric | Value |
|--------|-------|
| **Total Documents** | $total_docs |
| **Last Updated** | $(date +%Y-%m-%d' '%H:%M:%S) |

### By Dimension

| Dimension | Count |
|-----------|-------|
$(echo "$dim_stats" | awk '{print "| " $2 " " $3 " | " $1 " |"}')

### By Level

| Level | Count | Description |
|-------|-------|-------------|
$(echo "$level_stats" | awk '{printf "| %s | %s | %s |\n", $2, $1, ($2=="S"?"Expert":($2=="A"?"Advanced":($2=="B"?"Intermediate":"Basic")))}')

---

## 🗂️ Document Directory

EOF

    # Group by dimension
    declare -A dim_docs
    
    while IFS='|' read -r title dimension category level tags date path; do
        local dim_key="${dimension:-Other}"
        dim_docs["$dim_key"]+="|[$title](../$path)|$category|$level|$date|$path|\n"
    done < "$temp_file"
    
    # Output sorted by dimension
    for dim in "Formal Theory" "Language Design" "Engineering & Cloud Native" "Technology Stack" "Application Domains" "Examples" "Learning Paths" "Other"; do
        if [[ -n "${dim_docs[$dim]:-}" ]]; then
            echo "" >> "$output"
            echo "### $dim" >> "$output"
            echo "" >> "$output"
            echo "| Document | Category | Level | Date | Path |" >> "$output"
            echo "|----------|----------|-------|------|------|" >> "$output"
            echo -e "${dim_docs[$dim]}" | sort -t'|' -k4,4r >> "$output"
        fi
    done
    
    cat >> "$output" << 'EOF'

---

## 🔍 Quick Reference

### Document ID Prefixes

| Prefix | Dimension |
|--------|-----------|
| FT-* | Formal Theory |
| LD-* | Language Design |
| EC-* | Engineering & Cloud Native |
| TS-* | Technology Stack |
| AD-* | Application Domains |

### Level Definitions

| Level | Description | Target Audience |
|-------|-------------|-----------------|
| S | Expert/S-Level | Principal Engineers, Researchers |
| A | Advanced/A-Level | Senior Engineers |
| B | Intermediate/B-Level | Mid-level Engineers |
| C | Basic/C-Level | Junior Engineers |

---

*This index is automatically generated. Run `./scripts/generate-index.sh` to update.*
EOF

    log_success "Generated $output ($(wc -l < "$output") lines)"
}

# Update existing index files
update_indices_readme() {
    local output="$INDICES_DIR/README.md"
    
    log_info "Updating indices/README.md..."
    
    cat > "$output" << EOF
# Go Knowledge Base - Indices

> **Version**: Auto-generated
> **Last Updated**: $(date +%Y-%m-%d)

This directory contains auto-generated indices for the Go Knowledge Base.

## Available Indices

| Index | Description | File |
|-------|-------------|------|
| **By Topic** | Documents organized by topic/concept | [by-topic.md](./by-topic.md) |
| **By Tag** | Documents organized by tags | [by-tag.md](./by-tag.md) |
| **By Date** | Chronological listing of documents | [by-date.md](./by-date.md) |
| **Complete Map** | Full inventory with statistics | [complete-map.md](./complete-map.md) |
| **By Difficulty** | Learning paths by experience level | [by-difficulty.md](./by-difficulty.md) |
| **Cross Reference** | Topic relationships and pathways | [cross-reference.md](./cross-reference.md) |

## How to Update

Run the index generator script from the knowledge base root:

\`\`\`bash
./scripts/generate-index.sh
\`\`\`

Or from anywhere with the path:

\`\`\`bash
/path/to/scripts/generate-index.sh /path/to/go-knowledge-base
\`\`\`

## Index Details

### by-topic.md
Alphabetical index of topics with associated documents. Topics are extracted from:
- Document dimensions (维度)
- Document categories (分类)
- Document tags (标签)

### by-tag.md
Tag cloud style index with document counts per tag.

### by-date.md
Chronological view of document creation/updates, grouped by month.

### complete-map.md
Comprehensive document inventory including:
- Complete statistics
- All documents sorted by dimension
- Quick reference guides

---

*This README is automatically updated when indices are regenerated.*
EOF

    log_success "Updated $output"
}

# Generate statistics report
generate_stats() {
    local temp_file="$1"
    
    log_info "Generating statistics..."
    
    echo ""
    echo "========================================"
    echo "  Index Generation Statistics"
    echo "========================================"
    echo "Total Documents Scanned: $STATS_TOTAL"
    echo "Indices Generated:"
    echo "  - by-topic.md"
    echo "  - by-tag.md"
    echo "  - by-date.md"
    echo "  - complete-map.md"
    echo "========================================"
}

# Main function
main() {
    log_info "Starting index generation..."
    log_info "Knowledge Base Path: $KB_PATH"
    
    # Check if kb path exists
    if [[ ! -d "$KB_PATH" ]]; then
        log_error "Knowledge base path does not exist: $KB_PATH"
        exit 1
    fi
    
    # Ensure directories exist
    ensure_directories
    
    # Scan files
    local temp_file=$(scan_files)
    
    # Generate indices
    generate_by_topic "$temp_file"
    generate_by_tag "$temp_file"
    generate_by_date "$temp_file"
    generate_complete_map "$temp_file"
    update_indices_readme
    
    # Cleanup
    rm -f "$temp_file"
    
    # Show stats
    generate_stats "$temp_file"
    
    log_success "Index generation complete!"
    log_info "Indices location: $INDICES_DIR"
}

# Run main if not sourced
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi
