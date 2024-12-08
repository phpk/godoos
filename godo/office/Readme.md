# 📄 Gh0ffice (Office/PDF File Parser)

## Modifications
- 2024=12-08 godoos: add support for odt/epub/xml/rtf/md/txt/html/json files

This Go-based project provides a robust parser for various office document formats, including DOCX/DOC, PPTX/PPT, XLSX/XLS, and PDF. The parser extracts both content and metadata from these file types, allowing easy access to structured document data for further processing or analysis.

## 🛠 Features

- **Metadata Extraction**: Captures essential metadata such as title, author, keywords, and modification dates.
- **Content Parsing**: Supports extraction of text content from multiple file formats.
- **Extensible Architecture**: Easily add support for new file formats by implementing additional reader functions.

## 📂 Supported Formats

- **DOCX**: Extracts text content from Word documents.
- **PPTX**: Extracts text content from PowerPoint presentations.
- **XLSX**: Extracts data from Excel spreadsheets.
- **DOC**: Extracts text content from Legacy Word documents.
- **PPT**: Extracts text content from Legacy PowerPoint presentations.
- **XLS**: Extracts data from Legacy Excel spreadsheets.
- **PDF**: Extracts text content from PDF files (note that some complex PDFs may not be fully supported).

## 📖 Installation

To use this project, ensure you have Go installed on your system. Clone this repository and run the following command to install the dependencies:

```bash
go mod tidy
```

## 🚀 Usage

### Basic Usage

You can inspect a document and extract its content and metadata by calling the `inspectDocument` function with the file path as follows:

```go
doc, err := gh0ffice.InspectDocument("path/to/your/file.docx")
if err != nil {
    log.Fatalf("Error reading document: %s", err)
}
fmt.Printf("Title: %s\n", doc.Title)
fmt.Printf("Content: %s\n", doc.Content)
```

### Debugging

Set the `DEBUG` variable to `true` to enable logging for more verbose output during the parsing process:

```go
const DEBUG bool = true
```

## ⚠️ Limitations

- The PDF parsing may fail on certain complex or malformed documents.
- Only straightforward text extraction is performed; formatting and images are not considered.
- Compatibility tested primarily on major office file formats.

## 📝 License

This project is licensed under the Apache License, Version 2.0. See the [LICENSE](LICENSE) file for more details.

## 📬 Contributing

Contributions are welcome! Please feel free to create issues or submit pull requests for new features or bug fixes.

## 👥 Author

This project is maintained by the team and community of YT-Gh0st. Contributions and engagements are always welcome!

---

For any questions or suggestions, feel free to reach out. Happy parsing! 😊
