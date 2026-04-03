#!/usr/bin/env python3
"""
Go Knowledge Base - Automatic Index Generator

This script scans all markdown files and generates cross-referenced indices:
- indices/by-topic.md
- indices/by-tag.md
- indices/by-date.md
- indices/complete-map.md

Usage:
    python generate-index.py [knowledge_base_path]
    
Default path: parent of scripts directory
"""

import os
import re
import sys
from pathlib import Path
from datetime import datetime
from collections import defaultdict
from typing import List, Dict, Optional, NamedTuple


class DocumentMetadata(NamedTuple):
    """Represents metadata extracted from a markdown document."""
    title: str
    dimension: str
    category: str
    level: str
    tags: List[str]
    date: str
    path: str


class IndexGenerator:
    """Generates cross-referenced indices for the knowledge base."""
    
    # Dimension mapping from path prefixes
    DIMENSION_MAP = {
        "01-Formal-Theory": "Formal Theory",
        "02-Language-Design": "Language Design",
        "03-Engineering-CloudNative": "Engineering & Cloud Native",
        "04-Technology-Stack": "Technology Stack",
        "05-Application-Domains": "Application Domains",
        "examples": "Examples",
        "indices": "Indices",
        "learning-paths": "Learning Paths",
    }
    
    # Level thresholds based on file size
    LEVEL_THRESHOLDS = [
        (15000, "S"),
        (10000, "A"),
        (5000, "B"),
    ]
    
    def __init__(self, kb_path: str):
        self.kb_path = Path(kb_path).resolve()
        self.indices_dir = self.kb_path / "indices"
        self.documents: List[DocumentMetadata] = []
        
    def get_dimension_from_path(self, rel_path: str) -> str:
        """Infer dimension from file path."""
        parts = rel_path.replace("\\", "/").split("/")
        if parts:
            first_dir = parts[0]
            for prefix, dimension in self.DIMENSION_MAP.items():
                if first_dir.startswith(prefix) or first_dir == prefix:
                    return dimension
        return "Other"
    
    def get_level_from_size(self, file_size: int) -> str:
        """Infer level from file size."""
        for threshold, level in self.LEVEL_THRESHOLDS:
            if file_size >= threshold:
                return level
        return "C"
    
    def extract_metadata(self, file_path: Path) -> Optional[DocumentMetadata]:
        """Extract metadata from a markdown file."""
        try:
            rel_path = str(file_path.relative_to(self.kb_path)).replace("\\", "/")
            filename = file_path.name
            
            # Read first 30 lines for metadata
            with open(file_path, "r", encoding="utf-8", errors="ignore") as f:
                lines = []
                for i, line in enumerate(f):
                    if i >= 30:
                        break
                    lines.append(line)
            
            header = "".join(lines)
            
            # Extract title (first # line)
            title_match = re.search(r"^#\s+(.+)$", header, re.MULTILINE)
            title = title_match.group(1).strip()[:100] if title_match else filename
            
            # Extract dimension (维度) - only from quote blocks at the start
            dimension = ""
            dim_match = re.search(r"^>\s*\*\*维度\*\*[:：]\s*([^|\n]+)", header, re.MULTILINE)
            if dim_match:
                dimension = dim_match.group(1).strip()
            
            # Extract category (分类)
            category = ""
            cat_match = re.search(r"^>\s*\*\*分类\*\*[:：]\s*([^|\n]+)", header, re.MULTILINE)
            if cat_match:
                category = cat_match.group(1).strip()
            
            # Extract level (级别/难度)
            level = ""
            level_match = re.search(r"^>\s*\*\*级别\*\*[:：]\s*([SABC])", header, re.MULTILINE)
            if level_match:
                level = level_match.group(1)
            else:
                diff_match = re.search(r"^>\s*\*\*难度\*\*[:：]\s*(初级|中级|高级|Beginner|Intermediate|Advanced|Expert)", 
                                      header, re.MULTILINE)
                if diff_match:
                    diff_map = {
                        "初级": "B", "Beginner": "B",
                        "中级": "A", "Intermediate": "A",
                        "高级": "S", "Advanced": "S", "Expert": "S",
                    }
                    level = diff_map.get(diff_match.group(1), "B")
                elif "S-Level" in header or "S级" in header:
                    level = "S"
                elif "A-Level" in header or "A级" in header:
                    level = "A"
                elif "B-Level" in header or "B级" in header:
                    level = "B"
            
            # Extract tags - only from header metadata lines starting with >
            tags = []
            tag_line_match = re.search(r"^>\s*\*\*标签\*\*[:：]\s*(.+)$", header, re.MULTILINE)
            if tag_line_match:
                tag_text = tag_line_match.group(1)
                # Split by comma, space, or pipe
                for tag in re.split(r"[,\s|]+", tag_text):
                    tag = tag.strip().lstrip("#")
                    if tag and len(tag) > 1 and not tag.startswith(".."):
                        tags.append(tag.lower())
            
            # Extract date
            date = ""
            date_patterns = [
                r"^>\s*\*\*最后更新\*\*[:：]\s*(\d{4}-\d{2}-\d{2})",
                r"^>\s*\*\*完成日期\*\*[:：]\s*(\d{4}-\d{2}-\d{2})",
                r"^>\s*\*\*Created\*\*[:：]\s*(\d{4}-\d{2}-\d{2})",
            ]
            for pattern in date_patterns:
                date_match = re.search(pattern, header, re.MULTILINE)
                if date_match:
                    date = date_match.group(1)
                    break
            
            if not date:
                # Use file modification time
                mtime = file_path.stat().st_mtime
                date = datetime.fromtimestamp(mtime).strftime("%Y-%m-%d")
            
            # Infer dimension from path if not found
            if not dimension:
                dimension = self.get_dimension_from_path(rel_path)
            
            # Infer level from file size if not found
            if not level:
                file_size = file_path.stat().st_size
                level = self.get_level_from_size(file_size)
            
            return DocumentMetadata(
                title=title,
                dimension=dimension,
                category=category,
                level=level,
                tags=tags,
                date=date,
                path=rel_path
            )
            
        except Exception as e:
            print(f"  Warning: Error processing {file_path}: {e}")
            return None
    
    def scan_files(self) -> None:
        """Scan all markdown files in the knowledge base."""
        print(f"[INFO] Scanning markdown files in {self.kb_path}...")
        
        excluded_dirs = {"indices", "scripts", ".git"}
        
        count = 0
        for md_file in self.kb_path.rglob("*.md"):
            # Check if file is in excluded directory
            rel_parts = md_file.relative_to(self.kb_path).parts
            if any(part in excluded_dirs for part in rel_parts):
                continue
            
            metadata = self.extract_metadata(md_file)
            if metadata:
                self.documents.append(metadata)
                count += 1
                
                if count % 100 == 0:
                    print(f"\r[INFO] Scanned {count} files...", end="", flush=True)
        
        print(f"\r[SUCCESS] Scanned {count} markdown files")
    
    def escape_markdown(self, text: str) -> str:
        """Escape markdown special characters."""
        # Escape pipe characters in tables
        return text.replace("|", "\\|").replace("[", "\\[").replace("]", "\\]")
    
    def generate_by_topic(self) -> None:
        """Generate by-topic.md index."""
        output_path = self.indices_dir / "by-topic.md"
        print(f"[INFO] Generating by-topic.md...")
        
        # Build topic index
        topics = defaultdict(list)
        
        for doc in self.documents:
            # Add dimension as topic
            if doc.dimension:
                topics[doc.dimension].append(doc)
            
            # Add category as topic
            if doc.category:
                topics[doc.category].append(doc)
            
            # Add tags as topics
            for tag in doc.tags:
                if len(tag) > 1 and tag.isalnum() or "-" in tag:
                    topics[tag].append(doc)
        
        # Filter out invalid topics (too short or code snippets)
        valid_topics = {k: v for k, v in topics.items() 
                       if len(k) > 2 and not k.startswith(("func", "var", "const", "type", "if ", "for ", "// ", "{", "}"))}
        
        lines = [
            "# Go Knowledge Base - Topic-Based Cross Reference",
            "",
            f"> **Version**: Auto-generated",
            f"> **Last Updated**: {datetime.now().strftime('%Y-%m-%d')}",
            "> **Purpose**: Find documents by topic, concept, or technology",
            "",
            "---",
            "",
            "## 📑 Topics Overview",
            "",
        ]
        
        # Sort topics alphabetically
        sorted_topics = sorted(valid_topics.keys(), key=str.lower)
        
        # Group by first letter
        current_letter = ""
        for topic in sorted_topics:
            first_letter = topic[0].upper()
            if first_letter != current_letter:
                current_letter = first_letter
                lines.extend(["", f"### {current_letter}", ""])
            
            docs = valid_topics[topic][:5]  # Limit to 5 documents
            lines.append(f"- **{topic}**")
            
            for doc in docs:
                title = doc.title.replace("[", "\\[").replace("]", "\\]")
                lines.append(f"  - [{title}](../{doc.path})")
            
            if len(valid_topics[topic]) > 5:
                lines.append(f"  - *(and {len(valid_topics[topic]) - 5} more...)*")
            lines.append("")
        
        lines.extend([
            "---",
            "",
            "*This index is automatically generated. Run `./scripts/generate-index.py` to update.*",
        ])
        
        with open(output_path, "w", encoding="utf-8") as f:
            f.write("\n".join(lines))
        
        line_count = len(lines)
        print(f"[SUCCESS] Generated {output_path} ({line_count} lines)")
    
    def generate_by_tag(self) -> None:
        """Generate by-tag.md index."""
        output_path = self.indices_dir / "by-tag.md"
        print(f"[INFO] Generating by-tag.md...")
        
        # Collect tags
        tag_docs = defaultdict(list)
        
        for doc in self.documents:
            # Add tags
            for tag in doc.tags:
                clean_tag = tag.strip().lower()
                if clean_tag and len(clean_tag) > 1 and not clean_tag.startswith(("func", "var", "const")):
                    tag_docs[clean_tag].append(doc)
            
            # Add dimension as tag
            if doc.dimension:
                dim_tag = re.sub(r'[^a-z0-9]', '-', doc.dimension.lower()).strip('-')
                if dim_tag:
                    tag_docs[dim_tag].append(doc)
        
        lines = [
            "# Go Knowledge Base - Tag Index",
            "",
            f"> **Version**: Auto-generated",
            f"> **Last Updated**: {datetime.now().strftime('%Y-%m-%d')}",
            "> **Purpose**: Find documents by tags",
            "",
            "---",
            "",
            "## 🏷️ Tag Index",
            "",
        ]
        
        # Sort tags
        sorted_tags = sorted(tag_docs.keys(), key=str.lower)
        
        for tag in sorted_tags:
            count = len(tag_docs[tag])
            unique_docs = list({d.path: d for d in tag_docs[tag]}.values())
            
            lines.append(f"### `#{tag}` ({len(unique_docs)} documents)")
            lines.append("")
            
            for doc in sorted(unique_docs, key=lambda x: x.title):
                title = doc.title.replace("[", "\\[").replace("]", "\\]")
                lines.append(f"- [{title}](../{doc.path})")
            lines.append("")
        
        lines.extend([
            "---",
            "",
            "*This index is automatically generated. Run `./scripts/generate-index.py` to update.*",
        ])
        
        with open(output_path, "w", encoding="utf-8") as f:
            f.write("\n".join(lines))
        
        print(f"[SUCCESS] Generated {output_path} ({len(lines)} lines)")
    
    def generate_by_date(self) -> None:
        """Generate by-date.md index."""
        output_path = self.indices_dir / "by-date.md"
        print(f"[INFO] Generating by-date.md...")
        
        # Group by date
        date_docs = defaultdict(list)
        for doc in self.documents:
            if re.match(r"^\d{4}-\d{2}-\d{2}$", doc.date):
                date_docs[doc.date].append(doc)
        
        lines = [
            "# Go Knowledge Base - Chronological Index",
            "",
            f"> **Version**: Auto-generated",
            f"> **Last Updated**: {datetime.now().strftime('%Y-%m-%d')}",
            "> **Purpose**: Find documents by creation/update date",
            "",
            "---",
            "",
            "## 📅 Documents by Date",
            "",
        ]
        
        # Sort dates in reverse chronological order
        sorted_dates = sorted(date_docs.keys(), reverse=True)
        
        # Group by month
        current_month = ""
        for date in sorted_dates:
            month = date[:7]
            if month != current_month:
                current_month = month
                lines.extend(["", f"### {current_month}", ""])
                lines.append("| Date | Document | Dimension | Level |")
                lines.append("|------|----------|-----------|-------|")
            
            for doc in date_docs[date]:
                title = doc.title.replace("|", "\\|").replace("[", "\\[").replace("]", "\\]")
                lines.append(f"| {date} | [{title}](../{doc.path}) | {doc.dimension} | {doc.level} |")
        
        lines.extend(["", "## 📊 Monthly Statistics", ""])
        
        month_counts = defaultdict(int)
        for date in sorted_dates:
            month = date[:7]
            month_counts[month] += len(date_docs[date])
        
        for month in sorted(month_counts.keys(), reverse=True):
            lines.append(f"- **{month}**: {month_counts[month]} documents")
        
        lines.extend([
            "",
            "---",
            "",
            "*This index is automatically generated. Run `./scripts/generate-index.py` to update.*",
        ])
        
        with open(output_path, "w", encoding="utf-8") as f:
            f.write("\n".join(lines))
        
        print(f"[SUCCESS] Generated {output_path} ({len(lines)} lines)")
    
    def generate_complete_map(self) -> None:
        """Generate complete-map.md index."""
        output_path = self.indices_dir / "complete-map.md"
        print(f"[INFO] Generating complete-map.md...")
        
        total_docs = len(self.documents)
        dim_stats = defaultdict(int)
        level_stats = defaultdict(int)
        
        for doc in self.documents:
            dim = doc.dimension if doc.dimension else "Other"
            dim_stats[dim] += 1
            level_stats[doc.level if doc.level else "Unknown"] += 1
        
        lines = [
            "# Go Knowledge Base - Complete Document Map",
            "",
            f"> **Version**: Auto-generated",
            f"> **Last Updated**: {datetime.now().strftime('%Y-%m-%d')}",
            "> **Purpose**: Complete inventory of all knowledge base documents",
            "",
            "---",
            "",
            "## 📊 Statistics",
            "",
            "| Metric | Value |",
            "|--------|-------|",
            f"| **Total Documents** | {total_docs} |",
            f"| **Last Updated** | {datetime.now().strftime('%Y-%m-%d %H:%M:%S')} |",
            "",
            "### By Dimension",
            "",
            "| Dimension | Count |",
            "|-----------|-------|",
        ]
        
        for dim, count in sorted(dim_stats.items(), key=lambda x: -x[1]):
            lines.append(f"| {dim} | {count} |")
        
        lines.extend([
            "",
            "### By Level",
            "",
            "| Level | Count | Description |",
            "|-------|-------|-------------|",
        ])
        
        level_desc = {"S": "Expert", "A": "Advanced", "B": "Intermediate", "C": "Basic", "": "Unknown"}
        for level, count in sorted(level_stats.items(), key=lambda x: -x[1]):
            desc = level_desc.get(level, "Unknown")
            lines.append(f"| {level if level else '-'} | {count} | {desc} |")
        
        lines.extend([
            "",
            "---",
            "",
            "## 🗂️ Document Directory",
            "",
        ])
        
        # Group by dimension
        dim_order = [
            "Formal Theory", "Language Design", "Engineering & Cloud Native",
            "Technology Stack", "Application Domains", "Examples", 
            "Learning Paths", "Other"
        ]
        
        for dim in dim_order:
            dim_docs = [d for d in self.documents if (d.dimension or "Other") == dim]
            dim_docs.sort(key=lambda x: x.date, reverse=True)
            
            if dim_docs:
                lines.extend([
                    "",
                    f"### {dim} ({len(dim_docs)} documents)",
                    "",
                    "| Document | Category | Level | Date |",
                    "|----------|----------|-------|------|",
                ])
                
                for doc in dim_docs:
                    title = doc.title.replace("|", "\\|").replace("[", "\\[").replace("]", "\\]")
                    cat = doc.category if doc.category else "-"
                    lines.append(f"| [{title}](../{doc.path}) | {cat} | {doc.level} | {doc.date} |")
        
        lines.extend([
            "",
            "---",
            "",
            "## 🔍 Quick Reference",
            "",
            "### Document ID Prefixes",
            "",
            "| Prefix | Dimension |",
            "|--------|-----------|",
            "| FT-* | Formal Theory |",
            "| LD-* | Language Design |",
            "| EC-* | Engineering & Cloud Native |",
            "| TS-* | Technology Stack |",
            "| AD-* | Application Domains |",
            "",
            "### Level Definitions",
            "",
            "| Level | Description | Target Audience |",
            "|-------|-------------|-----------------|",
            "| S | Expert/S-Level | Principal Engineers, Researchers |",
            "| A | Advanced/A-Level | Senior Engineers |",
            "| B | Intermediate/B-Level | Mid-level Engineers |",
            "| C | Basic/C-Level | Junior Engineers |",
            "",
            "---",
            "",
            "*This index is automatically generated. Run `./scripts/generate-index.py` to update.*",
        ])
        
        with open(output_path, "w", encoding="utf-8") as f:
            f.write("\n".join(lines))
        
        print(f"[SUCCESS] Generated {output_path} ({len(lines)} lines)")
    
    def update_readme(self) -> None:
        """Update indices README.md."""
        output_path = self.indices_dir / "README.md"
        print(f"[INFO] Updating indices/README.md...")
        
        content = f"""# Go Knowledge Base - Indices

> **Version**: Auto-generated
> **Last Updated**: {datetime.now().strftime('%Y-%m-%d')}

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

```bash
# Using Python (recommended)
python scripts/generate-index.py

# Or with explicit path
python /path/to/scripts/generate-index.py /path/to/go-knowledge-base
```

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
"""
        
        with open(output_path, "w", encoding="utf-8") as f:
            f.write(content)
        
        print(f"[SUCCESS] Updated {output_path}")
    
    def run(self) -> None:
        """Run the index generation process."""
        print(f"[INFO] Starting index generation...")
        print(f"[INFO] Knowledge Base Path: {self.kb_path}")
        
        if not self.kb_path.exists():
            print(f"[ERROR] Knowledge base path does not exist: {self.kb_path}")
            sys.exit(1)
        
        # Ensure directories exist
        self.indices_dir.mkdir(exist_ok=True)
        print(f"[INFO] Indices directory: {self.indices_dir}")
        
        # Scan files
        self.scan_files()
        
        if not self.documents:
            print(f"[ERROR] No documents found to index!")
            sys.exit(1)
        
        # Generate indices
        self.generate_by_topic()
        self.generate_by_tag()
        self.generate_by_date()
        self.generate_complete_map()
        self.update_readme()
        
        # Show statistics
        print("")
        print("=" * 50)
        print("  Index Generation Complete!")
        print("=" * 50)
        print(f"Total Documents: {len(self.documents)}")
        print("Indices Generated:")
        print("  - by-topic.md")
        print("  - by-tag.md")
        print("  - by-date.md")
        print("  - complete-map.md")
        print("=" * 50)
        
        print(f"[SUCCESS] Index generation complete!")
        print(f"[INFO] Indices location: {self.indices_dir}")


def main():
    """Main entry point."""
    # Get knowledge base path from argument or use default
    if len(sys.argv) > 1:
        kb_path = sys.argv[1]
    else:
        # Default to parent of scripts directory
        script_dir = Path(__file__).parent
        kb_path = script_dir.parent
    
    generator = IndexGenerator(kb_path)
    generator.run()


if __name__ == "__main__":
    main()
