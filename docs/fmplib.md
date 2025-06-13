# FMP12 File Format Findings

(By Rasmus20B, 2024) - https://github.com/Rasmus20B/fmplib/blob/66245e5269275724bacfe1437fb1f73bc587a2f3/doc/fmp_format.md

## File Tree Structure

- Tables: [3].[16].[5]
- Relationships: [3].[17].[5]
- Fields: [table].[3].[5]
- Layout info: [4].[1].[7]
- Scripts: [17].[5].[script]
- value lists: [33].[5].[valuelist]

# Table Information

## Field type switches (Found at key 2 for field definition)

### Byte Index

### 0: Type of field (not data type, but simple, calculation, or summary)
- 0 = Simple field,
- 2 = Calculation field,
- 3 = Summary field,

### 1:  Simple Data-types.
- 1 = Text,
- 2 = Number,
- 3 = Date,
- 4 = Time,
- 5 = Timestamp,
- 6 = Container,

### 1: Summary Field Data-types.
- 1 = List of,
- 2 = Total of || Count of || Standard Deviation || Fraction of Total of,
- 5 = Average || Minimum || Maximum,

### 4:  Auto-Enter preset Options.
- 0 = Creation Date,
- 1 = Creation Time,
- 2 = Creation TimeStamp,
- 3 = Create Name,
- 4 = Creation Account Name,
- 5 = Modification Date,
- 6 = Modification Time,
- 7 = Modification TimeStamp,
- 8 = Modification Name,
- 9 = Modification Account Name,

### 7: Default language

2. Unicode,
3. Default,
16. Catalan,
17. Croatian,
18. Czech,
19. Danish,
20. Dutch,
21. English,
22. Finnish
23. Finnish (v\<\>w),
24. French,
25. German,
26. German (ä=a),
27. Greek,
28. Hungarian,
29. Icelandic,
30. Italian,
31. Japanese,
32. Norwegian,
33. Polish,
34. Portuguese,
35. Romanian,
36. Russian,
37. Slovak,
38. Slovenian,
39. Spanish (Modern),
40. Spanish,
41. Swedish,
42. Swedish (v\<\>w),
43. Turkish,
44. Ukrainian,
45. Chinese (Pinyin),
46. Chinese (Stroke),
47. Hebrew,
48. Hindi,
49. Arabic,
50. Estonian,
51. Lithuanian,
52. Latvian,
53. Serbian (Latin),
54. Farsi,
55. Bulgarian,
56. Vietnamese,
57. Thai,
58. Greek (Mixed),
59. Bengali,
60. Telugu,
61. Marathi,
62. Tamil,
63. Gujarati,
64. Kannada,
65. Malayalam
67. Panjabi,
76. Korean,


### 8:
- 64 = Don't automatically create index,
- 128 = ALways index this field,

### 9:
- 0 = regular storage,
- 1 = Global Field,
- 8 = Calculation field,
- 10 = Unstored Calculation,

### 10:
- 1: Prohibit modification of value during data-entry
- 2: Serial Number **On Commit**,
- 4: Set in conjunction with idx 11: 128 to signify lookup,

### 11:
- 1 = Options flag from idx 4,
- 2 = Serial Number **On Creation**,
- 4 = Data textbox,
- 8 = Auto-Enter Calculation (**does not** replace existing value),
- 16 = Value from last visited record,
- 32 = Evaluate Calculation even if all referenced fields are empty,
- 128 = with idx 10 = 4 it is a lookup that is active, otherwise inactive but data saved elsewhere,
- 136 = Auto-Enter Calculation (**does** replace existing value),

### 14:
- 0 = Only validate during data entry,
- 1 = Member of value list,
- 2 = maximum number of characters,
- 4 = always validate,
- 16 = strict data-type: Numeric only,
- 32 = strict data-type: 4 digit year,
- 64 = strict data-type: Time of Day,

### 15:
- 0 = User can override,
- 1 = Validated by calculation,
- 4 = User cannot override,
- 8 = Required value,
- 16 = Unique Value,
- 32 = Existing Value,
- 64 = within a range of values,
- 128 = Display a validation error message,

### 25:
- byte 25 simply states how many repetitions the field has

# Relationships

## Relationship Structure

### [3].[17].[5].[0]

- Each relationship is read sequentially at the same path, reusing the same key for references.
- Graph mapping data is stored at path **[3].[17].[5].[0].[251]**.

### Keys used in each relationship

- (2) => 35 bytes that specify metadata about the current table occurence.
    - Byte 7 specificies table that occurence is based on.
- (16) => Name of the table occurence.
- (216) => gimme some time
- [3].[17].[5].[0].[251] => Simple Data. typically 5 Bytes.
- (252) => ???

# Calculation Engine

Calculations are stored in a kind of bytecode, with basic operators ('+', '-', etc) being encoded as ints.

## Operators

- '+' :: 0x25
- '-' :: 0x26
- '\*' :: 0x27
- '/' :: 0x28
- '&' :: 0x50

## How to decode numbers
Numbers start with a 0x10. The 9th byte will be the first byte of the number

## How to decode variables
Variables start with '0x1a', followed by the size of the variable name string.

# Scripts

## Scripting Structure

- 4 will usually denote the actual top-level script code.

### Script code
- Each step is stored as a 24 byte subarray, most commonly starting with '2, 1'.
- Bytes 3 and 4 are used to index the script step. This 'index' can be used in the script step 'data' directory specified below.
- **Important**: When script runs into space constraints, simple key ref does not suffice. Segments of the array are stored at **Path** [17].[5].[script].[4], rather than key-value.


### [17].[1].[7].[script]
- This path contains a small amount of metadata for the scripts, most notably their names located at key 16.

### [17].[5].[script]::4
- This path stores the instructions from the script, some metadata, as well as
the file value to append to the path above.

- Script attributes should be parsed in order of appearance in file format. If the script content appears first, store it's id with the content, with a blank name field and other options. Then when we find the name in the metadata directory, just look for the corresponding script and fill in the name.

### [17].[5].[script].[4]
- Used when script exceeds size limitations. replaces key of 4 in script directory ([17].[5].[script]).

### [17].[5].[script].[5]
- Store information about each script step. Each script step being a number after the 5, and treated as it's own directory.

### List of Instructions
#### 1. Perform Script
#### 3. Save a Copy as XML
#### 4. Go to Next Field
#### 5. Go to Previous Field
#### 6. Go to Layout
    - bytes 7, 8, and 9: Specify layout. If "original layout" these are all zero.
#### 7. New Record/Request
#### 8. Duplicate Record/Request
#### 9. Delete Record/Request
#### 10. Delete All Records
#### 11. Insert From Index
#### 12. Insert From Last Visited
#### 13. Insert Current Date
#### 14. Insert Current Time
#### 16. Go to Record/Request/Page
#### 17. Go to Field
#### 18. Check Selection
#### 19. Check Record
#### 20. Check Found Set
#### 21. Unsort Records
#### 22. Enter Find Mode
#### 23. Show All Records
#### 24. Modify Last Find
#### 25. Omit Record
#### 26. Omit Multiple Records
#### 27. Show Ommited Only
#### 28. Perform Find
#### 29. Show/Hide Toolbars
#### 30. View As
#### 31. Adjust Window
#### 32. Open Help
#### 33. Open File
#### 34. Close File
#### 35. Import Records
#### 36. Export Records
#### 37. Save a Copy as
#### 38. Open Manage Database
#### 39. Sort Records
#### 40. Relookup Field Contents
#### 41. Enter Preview Mode
#### 42. Print Setup
#### 43. Print
#### 44. Exit Application
#### 45. Undo/Redo
#### 46. Cut
#### 47. Copy
#### 48. Paste
#### 49. Clear
#### 50. Select All
#### 51. Revert Record/Request
#### 55. Enter Browser Mode
#### 56. Insert Picture
#### 57. Send Event
#### 60. Insert Current User Name
#### 61. Insert Text
#### 62. Pause/Resume Script
#### 63. Send Mail
#### 64. Send DDE Execute
#### 65. Dial Phone
#### 66. Speak
#### 67. Perform AppleScript
#### 68. If
#### 69. Else
#### 70. End If
#### 71. Loop
#### 72. Exit Loop If
#### 73. End Loop
#### 74. Go to Related Record
#### 75. Commit Records/Requests
#### 76. Set Field
#### 77. Insert Calculated Result
#### 79. Freeze Window
#### 80. Refresh Window
#### 81. Scroll Window
#### 82. New File
#### 83. Change Password
#### 84. Set Multi-User
#### 85. Allow User abort. 26th byte option
    - 26: 1 = Off, 3 = On.
#### 86. Set Error Capture
#### 87. Show Custom Dialog
#### 88. Open Script Workspace
#### 89. Blank Line/Comment
#### 90. Halt Script
#### 91. Replace Field Contents
#### 92. Show/Hide Text Ruler
#### 93. Beep
#### 94. Set Use System Formats
#### 95. Recover File
#### 96. Save a Copy as Add-on Package
#### 97. Set Zoom Level
#### 98. Copy All Records/Requests
#### 99. Go to Portal Row
#### 101. Copy Record/Request
#### 102. Fluch Cache to Disk
#### 103. Exit Script
#### 104. Delete Portal Row
#### 105. Open Preferences
#### 106. Correct Word
#### 107. Spelling Options
#### 108. Select Dictionaries
#### 109. Edit User Dictionary
#### 111. Open URL
#### 112. Open Manage Value Lists
#### 113. Open Sharing
#### 114. Open File Options
#### 115. Allow Formatting Bar
#### 116. Set Next Serial Value
#### 117. Execute SQL
#### 118. Open Hosts
#### 119. Move/Resize Window
#### 120. Arrange All Windows
#### 121. Close Window
#### 122. New Window
#### 123. Select Window
#### 124. Set Window Title
#### 125. Else If
#### 126. Constrain Found Set
#### 127. Extend Found Set
#### 128. Perform Find/Replace
#### 129. Open Find/Replace
#### 130. Set Selection
#### 131. Insert File
#### 132. Export Field Contents
#### 133. Open Record Request
#### 134. Add Account
#### 135. Delete Account
#### 136. Reset Account Password
#### 137. Enable Account
#### 138. Relogin
#### 139. Convert File
#### 140. Open Manage Data Sources
#### 141. Set Variable
    - Stores variable name to be set @ [17].[5].[scriptnumber].[5].[instructionnumber].[128]::1
    - stores new value calculation @ [17].[5].[scriptnumber].[5].[instructionnumber].[129].[5]::5
#### 142. Install Menu Set
#### 143. Save Records as Excel
#### 144. Save Records as PDF
#### 145. Go to Object
#### 147. Set Field by Name
#### 146. Set Web Viewer
#### 148. Install OnTimer Script
#### 149. Open Edit Saved Finds
#### 150. Perform Quick Find
#### 151. Open Manage Layouts
#### 152. Save Records as Snapshot Link
#### 154. Sort Records by Field
#### 155. Find Matching Records
#### 156. Manage Containers
#### 157. Install Plugin File
#### 158. Insert PDF
#### 159. Insert Audio/Video
#### 160. Insert from URL
#### 161. Insert From Device
#### 164. Perform Script on Server
#### 165. Open Manage Themes
#### 166. Show/Hide Menubar
#### 167. Refresh Object
#### 168. Set Layout Object Animation
#### 169. Close Popover
#### 172. Open Upload to Host
#### 174. Enable Touch Keyboard
#### 175. Perform JavaScript in Web Viewer
#### 177. AVPlayer Play
#### 178. AVPlayer Set Playback State
#### 179. AvPlayer Set Options
#### 180. Refresh Portal
#### 181. Get Folder Path
#### 182. Truncate Table
#### 183. Open Favorites
#### 185. Configure Region Monitor Script
#### 187. Configure Local Notification
#### 188. Get File Exists
#### 189. Get File Size
#### 190. Create Data File
#### 191. Open Data File
#### 192. Write to Data File
#### 193. Read from Data File
#### 194. Get Data File Position
#### 195. Set Data File Position
#### 196. Close Data File
#### 197. Delete File
#### 199. Rename File
#### 200. Set Error Logging
#### 201. Configure NFC Reading
#### 202. Configure Machine Learning Model
#### 203. Execute FileMaker Data API
#### 205. Open Transaction
#### 206. Commit Transaction
#### 207. Revert Transaction
#### 208. Set Session Identifier
#### 209. Set Dictionary
#### 210. Perform Script on Server with Callback
#### 211. Trigger Claris Connect Flow

### [17].[5].[script].[5] - The Instruction Directory
- The "data" for each script step is located in this folder.
