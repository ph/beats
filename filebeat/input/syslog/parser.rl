// Code generated by ragel DO NOT EDIT.
package syslog

%%{
    machine syslog;
    write data;
    variable p p;
    variable pe pe;
}%%

// syslog
//<34>Oct 11 22:14:15 wopr su: 'su root' failed for foobar
//<13>Feb  5 17:32:18 10.0.0.99 Use the quad dmg.
func Parse(data []byte, event *event) {
    var p, cs int
    pe := len(data)
    tok := 0
    eof := len(data)
    %%{
      action tok {
        tok = p
      }

      action priority {
        event.SetPriority(data[tok:p])
      }

      action message {
        event.SetMessage(data[tok:p])
      }

      action month {
        event.SetMonth(data[tok:p])
      }

      action year{
        event.SetYear(data[tok:p])
      }

      action month_numeric {
        event.SetMonthNumeric(data[tok:p])
      }

      action day {
        event.SetDay(data[tok:p])
      }

      action hour {
        event.SetHour(data[tok:p])
      }

      action minute {
        event.SetMinute(data[tok:p])
      }

      action second {
        event.SetSecond(data[tok:p])
      }

      action nanosecond{
        event.SetNanosecond(data[tok:p])
      }

      action hostname {
        event.SetHostname(data[tok:p])
      }

      action program {
        event.SetProgram(data[tok:p])
      }

      action pid {
        event.SetPid(data[tok:p])
      }

      action timezone { }

      include syslog_rfc3164 "syslog_rfc3164.rl";

      write init;
      write exec;
    }%%
}
