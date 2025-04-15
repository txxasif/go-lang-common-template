# Creating Database Diagrams in draw.io

## Your Database Schema in draw.io

To create your database diagram in draw.io, follow these steps:

1. Open draw.io (either desktop or web version)
2. Create a new diagram
3. Select "Entity Relation" from the template gallery

## Step-by-Step Instructions

### 1. Creating the Users Table

1. From the left panel, drag and drop a "Table" shape onto the canvas
2. Double-click the table to edit it
3. Add the following structure:

```
Users
--------
id (PK)
username
role
created_at
```

### 2. Creating the Follows Table

1. Add another "Table" shape
2. Double-click to edit and add:

```
Follows
--------
following_user_id (FK)
followed_user_id (FK)
created_at
```

### 3. Adding Relationships

1. From the left panel, select the "Relationship" connector
2. Draw a line from the `id` field in Users to `following_user_id` in Follows
3. Draw another line from the `id` field in Users to `followed_user_id` in Follows

### 4. Styling the Diagram

1. Right-click on tables to:
   - Change colors
   - Adjust size
   - Modify text formatting
2. Use the format panel to:
   - Add shadows
   - Change line styles
   - Adjust spacing

## Visual Representation

Your diagram should look like this:

```
+-------------+       +-------------+
|    Users    |       |   Follows   |
+-------------+       +-------------+
| id (PK)     |<----->| following_  |
| username    |       | user_id (FK)|
| role        |       | followed_   |
| created_at  |<----->| user_id (FK)|
+-------------+       | created_at  |
                     +-------------+
```

## Tips for Better Diagrams

1. **Table Styling**

   - Use consistent colors for primary and foreign keys
   - Add a different background color for the header row
   - Use bold text for primary keys

2. **Relationship Lines**

   - Use solid lines for required relationships
   - Add crow's foot notation for cardinality
   - Label relationships clearly

3. **Layout**

   - Arrange tables to minimize line crossings
   - Group related tables together
   - Use grid alignment for a clean look

4. **Documentation**
   - Add a title to your diagram
   - Include a legend for symbols
   - Add notes for complex relationships

## Export Options

1. **File Formats**

   - PNG (for documentation)
   - PDF (for printing)
   - SVG (for web use)
   - XML (for editing later)

2. **Sharing**
   - Export to cloud storage
   - Share via link
   - Embed in documentation

## Common Issues and Solutions

1. **Overlapping Tables**

   - Use the "Arrange" menu to space tables evenly
   - Enable grid snapping for alignment
   - Use the "Layout" feature to auto-arrange

2. **Complex Relationships**

   - Use different line styles for different relationship types
   - Add intermediate tables for many-to-many relationships
   - Use notes to explain complex constraints

3. **Readability**
   - Increase font size for better visibility
   - Use contrasting colors
   - Add white space between tables
