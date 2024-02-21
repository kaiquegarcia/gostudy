# gostudy

A simple study-plan maker made in Golang.

## How to use

1. You must have [Golang v1.20+](https://go.dev/) installed on your machine. If you don't have, access its page, download and install it;
2. Clone this repository;
3. Use your folder explorer (Windows Explorer, Finder, whatever) to access the root path of the cloned repository;
4. Prepare your datasheets:
    1. hour grade:
        * copy the [templates_hour_grade.csv](./templates_hour_grade.csv) to a new file `hour_grade.csv`;
        * write all the time intervals you of your study routine to `hour_grade.csv` following the format `hh:mm-hh:mm`. The first `hh:mm` is the start time and the last is the limit;
        * each line of the table is a day of week;
        * from the second column to the end, you can put all fragmented intervals you want to study;
        * don't break your intervals with "gaps", as the system can automatically add gaps during the plan-making;
        * you can put as many intervals as you want, but don't leave blank spaces from an interval to another (example: `15:00-16:00,,19:00-20:00` ~ the `,,` will make it ignore the `19:00-20:00` interval).
    2. disciplines list:
        * copy the [template_disciplines.csv](./template_disciplines.csv) to a new file `disciplines.csv`;
        * write all disciplines you'll study on this plan to `disciplines.csv`;
        * remember to write the `filename` correctly, so the system can find the discipline's content file;
        * `daily limit` means how many hours/minutes/seconds you accept to have content from this disciplines `per day`;
        * `content gap` means how many hours/minutes/seconds you want to append before each content of this discipline, except for the first content of the time interval;
        * `subject gap` means how many hours/minutes/seconds you want to append before each subject change for this discipline, except for the first content of the time interval.
    3. disciplines contents:
        * based on the `filenames` you written on `disciplines.csv`, copy [template_{discipline_file}.csv](./template_{discipline_file}.csv) for each `filename` present on `disciplines.csv`;
        * write all content you will study there in order of study;
        * the `Subject` will be the key to group the contents by subject (to know when to use discipline's `subject gap`);
        * the `Duration` is also a key for the plan-maker to properly place the content on the intervals. **If you put an unplayable duration, the plan-maker will return error after exceed attempts of putting the content on the plan**. For example, if you only study 1 hour per day but have a content with 2 hours of duration, it won't be reachable, resulting on error.
5. Open a terminal on the root path of the cloned repository;
6. Run `go run .` and follow the software instructions!

## Changing the initial date

The initial date of the plan is, by default, the same current day of next week (base on your machine's datetime).

You can change the initial date by sending it on the run command.

For example: `go run . 2024-02-21`, which should have the start date as `2024-02-21`.

If the start date doesn't have any time interval on the hour grade, it will get the very next date with available time interval.