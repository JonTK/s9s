# Export & Reporting Guide

Export S9S data in supported formats for analysis, reporting, and integration with other tools and workflows.

## Quick Export

### Basic Export Commands

```bash
# Export current view to CSV
:export csv

# Export selected jobs to JSON
/user:alice state:COMPLETED
:export --selected json

# Export with custom filename
:export csv --output my-jobs.csv

# Export to Markdown
:export md --output report.md

# Export to plain text
:export txt --output jobs.txt
```

## Supported Formats

S9S currently supports the following export formats:

### CSV (Comma-Separated Values)
- Best for: Spreadsheet analysis, data processing
- Features: Headers, UTF-8 encoding
- Extension: `.csv`

**Example**:
```bash
:export csv --output jobs.csv
```

### JSON (JavaScript Object Notation)
- Best for: API integration, web applications, programmatic processing
- Features: Structured data, nested information
- Extension: `.json`

**Example**:
```bash
:export json --output jobs.json
```

### Markdown
- Best for: Documentation, version control, readable reports
- Features: Tables, formatting, human-readable
- Extension: `.md`

**Example**:
```bash
:export md --output report.md
```

### Plain Text
- Best for: Simple logs, basic reporting, console output
- Features: Simple formatting, universally readable
- Extension: `.txt`

**Example**:
```bash
:export txt --output jobs.txt
```

## Export Options

### Field Selection

Choose specific fields to export:

```bash
# Export specific job fields
:export csv --fields=JobID,User,State,Runtime

# Export all available fields
:export json --fields=all

# Export minimal fields
:export csv --fields=minimal
```

### Time Range Filters

```bash
# Export jobs from specific time period
:export csv --time-range="2023-12-01..2023-12-31"

# Export recent data
:export json --time-range="last-7d"

# Export with relative time
:export csv --submitted=">1h" --completed="<24h"
```

### Data Filtering

```bash
# Export with filters
:export csv --filter="user:alice state:COMPLETED"

# Complex filtering
:export json --filter="partition:gpu nodes:>4 runtime:>2h"

# Multiple conditions
:export csv --user=alice,bob --state=RUNNING,COMPLETED
```

## Export Destinations

### Local Files

```bash
# Export to specific directory
:export csv --output="/data/exports/jobs.csv"

# Export with timestamp
:export json --output="jobs-{timestamp}.json"

# Export to user directory
:export md --output="~/reports/cluster-report.md"
```

## Export Configuration

### Default Settings

```yaml
# ~/.s9s/config.yaml
export:
  # Default format
  defaultFormat: csv

  # Default output directory
  outputDir: ~/s9s-exports

  # Include headers in CSV
  includeHeaders: true

  # Date format in filenames
  dateFormat: "2006-01-02"

  # Field formatting
  timeFormat: RFC3339
  durationFormat: seconds

  # Limits
  maxRecords: 1000000
  maxFileSize: 100MB
```

### Format-Specific Settings

```yaml
export:
  formats:
    csv:
      delimiter: ","
      quote: '"'
      encoding: utf-8
      lineEnding: unix

    json:
      indent: 2
      sortKeys: true
      includeSchema: false

    markdown:
      tableFormat: github
      includeHeaders: true

    txt:
      lineFormat: simple
      separator: " | "
```

## Security and Privacy

### Data Sanitization

```bash
# Remove sensitive information
:export csv --sanitize --fields=JobID,State,Runtime

# Filter sensitive partitions
:export csv --exclude-partitions=confidential,private
```

### Access Control

```yaml
export:
  security:
    requirePermission: true
    allowedFormats: [csv, json, md, txt]
    maxRecordsPerUser: 10000
    auditExports: true
    restrictFields: [script_path, environment]
```

## Best Practices

### Performance Optimization

1. **Use filters** to limit data volume
2. **Export incrementally** for large historical data
3. **Choose appropriate formats** (JSON for structured data, CSV for spreadsheets, Markdown for reports)
4. **Limit field selection** to only needed columns

### Data Management

1. **Version your exports** with timestamps
2. **Archive old exports** regularly
3. **Document export schemas** for consistency
4. **Validate exported data** before use
5. **Monitor export jobs** for failures

## Troubleshooting

### Common Issues

**Export fails with "Permission denied"**:
- Check file/directory permissions
- Verify export destination accessibility
- Ensure sufficient disk space

**Large exports timeout**:
- Use smaller time ranges
- Export in batches
- Apply filters to reduce data volume

**Invalid date formats**:
- Check date format configuration
- Verify timezone settings
- Use ISO 8601 format for compatibility

### Debug Mode

```bash
# Enable export debugging
:config set export.debug true

# Verbose export logging
:export csv --debug --verbose

# Dry run export
:export json --dry-run --output=test.json
```

## Planned Features

The following features are planned for future releases:

### Additional Formats (Planned)
- **Excel (.xlsx)** - For business reports and formatted presentations
- **PDF** - For formatted reports and documentation
- **Parquet** - For big data analysis and columnar storage
- **HTML** - For web reports and interactive content

### Advanced Reporting (Planned)
- Job summary reports with statistics and charts
- User activity reports with resource consumption metrics
- Resource utilization reports with efficiency analysis
- Node health and maintenance reports
- Custom report templates

### Automated Exports (Planned)
- Scheduled exports (daily, weekly, monthly)
- Event-triggered exports
- Export automation with conditions

### Cloud Storage Integration (Planned)
- AWS S3
- Google Cloud Storage
- Azure Blob Storage

### Database Integration (Planned)
- PostgreSQL export
- InfluxDB export
- Elasticsearch export

### API Endpoints (Planned)
- Webhook POST
- Apache Kafka streaming
- Prometheus metrics

### Data Visualization (Planned)
- Built-in charts (pie, line, bar)
- Integration with Tableau, Power BI, Grafana

## Next Steps

- Learn [Batch Operations](./batch-operations.md) for bulk exports
- Explore [Node Operations](./node-operations.md) for node data analysis
- Explore [Advanced Filtering](../filtering.md) to refine export data
