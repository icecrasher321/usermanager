# usermanager
A repository for handling various user data.


Problem Statement
----------------

This is a user management system. It will read from some source the incoming user entries and logs it in some sort of storage.


User fields
-----------
    `username        String`
    `FirstName       String`
    `LastName        String`
    `Age             Integer`
    `Mobile No.[s]   Integer`
    `Email Id[s]     String`

Notes
-----
    *   username must be unique.
    *   A user can have multiple mobile numbers and emails
    *   There should be validations on Age, Mobile no. and EmailIds

Stage - I
---------
    *   Application will take entries from commandline arguments/flags.
    *   Application will have these 3 functionality.
        -   Create
        -   Update
        -   Delete
    *   Application will write into a file.


Stage - II
----------
    *   Application will accept entries from file. Read the file and do operations.
            <Description would be filled later.>
    *   There should be no duplicate user entry.
    *   Application should be able to give user details on query. "No entry found" in case of no data.


Stage - III
-----------
    *   Application should be able to read big files. Size of files could be more than 10000.
    *   Application should be able to minimize functionality time.
            Hint -  Multithreaded environment could be used for that.
    *   Application should have LRU Cache (Least Recently Used) inside it.
            <We will discuss about it later.>


Stage - IV [Optional]
---------------------
    *   Application should have a functionality for uploading a file and return with status. [Use of GIN]
