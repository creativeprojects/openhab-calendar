# openhab-calendar

Simple script to load information from an iCalendar into OpenHAB.

The iCal binding (in OpenHAB v2) didn't work for me because:
- it can't use a Digest authentication
- it's a bit tedious to setup rules in Xtend
- it's much easier to debug on my local machine

Example configuration file:
```json
{
    "rules": [
        {
            "priority": 10,
            "name": "Override",
            "calendar": {
                "url": "http://calendars/override"
            },
            "result": "OVERRIDE"
        },
        {
            "priority": 20,
            "name": "Away",
            "calendar": {
                "url": "http://calendars/away"
            },
            "result": "AWAY"
        },
        {
            "priority": 30,
            "name": "DayOff",
            "weekdays": [
                "Mon",
                "Tue",
                "Wed",
                "Thu",
                "Fri"
            ],
            "calendar": {
                "url": "http://calendars/dayoff"
            },
            "result": "DAYOFF"
        },
        {
            "priority": 40,
            "name": "Office",
            "weekdays": [
                "Mon",
                "Tue",
                "Wed",
                "Thu",
                "Fri"
            ],
            "calendar": {
                "url": "http://calendars/office"
            },
            "result": "OFFICE"
        },
        {
            "priority": 50,
            "name": "WorkFromHome",
            "weekdays": [
                "Mon",
                "Tue",
                "Wed",
                "Thu",
                "Fri"
            ],
            "result": "WFH"
        },
        {
            "priority": 60,
            "name": "Weekend",
            "weekdays": [
                "Sat",
                "Sun"
            ],
            "result": "DAYOFF"
        }
    ],
    "post-rules": [
        {
            "priority": 10,
            "name": "First day of week-end",
            "when": { "is": "DAYOFF"},
            "previous": { "not": "DAYOFF" },
            "next": { "is": "DAYOFF" },
            "result": "JOBBIES"
        }
    ],
    "default": {
        "name": "Unknown",
        "result": "ERROR"
    },
    "servers": {
        "http": {
            "listen": "http://:6060"
        }
    },
    "authentication": [
        {
            "url": "http://calendars/",
            "username": "username",
            "password": "password"
        }
    ]
}
```