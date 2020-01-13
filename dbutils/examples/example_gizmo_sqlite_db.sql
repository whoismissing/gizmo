BEGIN TRANSACTION;
CREATE TABLE IF NOT EXISTS "Service" (
	"Name"	TEXT NOT NULL UNIQUE,
	"TeamID"	INTEGER,
	"ServiceID"	INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
	"HostIP"	TEXT,
	"NumberOfMissedChecks"	INTEGER,
	"NumberOfChecks"	INTEGER,
	"User"	TEXT,
	"Password"	TEXT,
	"Domain"	TEXT
);
CREATE TABLE IF NOT EXISTS "Status" (
	"ServiceID"	INTEGER,
	"StatusID"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
	"Time"	INTEGER,
	"State"	TEXT
);
CREATE TABLE IF NOT EXISTS "Team" (
	"GameID"	INTEGER,
	"TeamID"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
	"TotalMissedChecks"	INTEGER,
	"TotalChecks"	INTEGER
);
CREATE TABLE IF NOT EXISTS "Game" (
	"GameStartTime"	INTEGER,
	"CurrentGameTime"	INTEGER,
	"GameID"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE
);
INSERT INTO "Service" ("Name","TeamID","ServiceID","HostIP","NumberOfMissedChecks","NumberOfChecks","User","Password","Domain") VALUES ('www',1,1,'127.0.0.1',0,0,NULL,NULL,'localhost');
INSERT INTO "Service" ("Name","TeamID","ServiceID","HostIP","NumberOfMissedChecks","NumberOfChecks","User","Password","Domain") VALUES ('ssh1',1,2,'127.0.0.1',0,0,'foo','bar',NULL);
INSERT INTO "Status" ("ServiceID","StatusID","Time","State") VALUES (1,1,1566518090,'UP');
INSERT INTO "Team" ("GameID","TeamID","TotalMissedChecks","TotalChecks") VALUES (1,1,0,0);
INSERT INTO "Game" ("GameStartTime","CurrentGameTime","GameID") VALUES (1566517722,1566517722,1);
COMMIT;
