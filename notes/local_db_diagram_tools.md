# Local Database Diagram Tools for Ubuntu

## 1. DBeaver

A powerful database tool that includes ER diagram functionality.

### Installation

```bash
# Add DBeaver repository
sudo add-apt-repository ppa:serge-rider/dbeaver-ce
sudo apt update
sudo apt install dbeaver-ce
```

### Features

- Free and open-source
- Supports multiple database systems
- Auto-generates diagrams from existing databases
- Interactive diagram editing
- Export to various formats (PNG, SVG, PDF)
- Reverse engineering from existing databases

## 2. MySQL Workbench

Excellent for MySQL/MariaDB database design.

### Installation

```bash
sudo apt update
sudo apt install mysql-workbench
```

### Features

- Free and open-source
- Visual database design
- Forward and reverse engineering
- Schema synchronization
- SQL development
- Database administration

## 3. pgModeler

Specialized tool for PostgreSQL database modeling.

### Installation

```bash
# Add repository
sudo sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list'
wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo apt-key add -
sudo apt update
sudo apt install pgmodeler
```

### Features

- Free and open-source
- PostgreSQL-specific features
- Visual modeling
- SQL code generation
- Reverse engineering
- Export to various formats

## 4. Draw.io (Desktop)

Popular diagramming tool with database diagram support.

### Installation

```bash
# Download the AppImage
wget https://github.com/jgraph/drawio-desktop/releases/download/v20.8.16/drawio-amd64-20.8.16.AppImage
chmod +x drawio-amd64-20.8.16.AppImage
./drawio-amd64-20.8.16.AppImage
```

### Features

- Free and open-source
- Extensive template library
- Database diagram templates
- Export to multiple formats
- Cloud sync option
- Collaborative features

## 5. Dia

Simple diagramming tool with database support.

### Installation

```bash
sudo apt update
sudo apt install dia
```

### Features

- Free and open-source
- Lightweight
- Basic database diagram support
- Export to various formats
- Simple interface

## 6. SchemaSpy

Command-line tool for database schema visualization.

### Installation

```bash
sudo apt update
sudo apt install schemaspy
```

### Features

- Free and open-source
- Command-line based
- Generates HTML documentation
- Supports multiple databases
- Detailed relationship analysis

## Comparison Table

| Tool            | Type | Database Support | Ease of Use | Export Formats | Reverse Engineering |
| --------------- | ---- | ---------------- | ----------- | -------------- | ------------------- |
| DBeaver         | GUI  | Multiple         | High        | Multiple       | Yes                 |
| MySQL Workbench | GUI  | MySQL/MariaDB    | High        | Multiple       | Yes                 |
| pgModeler       | GUI  | PostgreSQL       | Medium      | Multiple       | Yes                 |
| Draw.io         | GUI  | Generic          | High        | Multiple       | No                  |
| Dia             | GUI  | Generic          | Medium      | Multiple       | No                  |
| SchemaSpy       | CLI  | Multiple         | Low         | HTML           | Yes                 |

## Recommendations

1. **For General Database Design**

   - DBeaver: Best all-around tool with support for multiple databases
   - Draw.io: Great for quick diagrams and collaboration

2. **For Specific Databases**

   - MySQL Workbench: Best for MySQL/MariaDB
   - pgModeler: Best for PostgreSQL

3. **For Documentation**

   - SchemaSpy: Best for generating detailed documentation
   - DBeaver: Good for both design and documentation

4. **For Simple Diagrams**
   - Dia: Lightweight and simple
   - Draw.io: User-friendly with templates

## Best Practices

1. **Version Control**

   - Store diagram files in version control
   - Use text-based formats when possible
   - Include documentation with diagrams

2. **Collaboration**

   - Use tools that support collaborative editing
   - Share diagrams in common formats
   - Maintain consistent naming conventions

3. **Documentation**

   - Include table descriptions
   - Document relationships
   - Add notes for complex structures

4. **Maintenance**
   - Keep diagrams up to date
   - Review regularly
   - Update documentation with changes
