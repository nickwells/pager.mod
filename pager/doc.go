/*
Package pager provides types and methods for paging output. It takes a value
which satisfies the SetWriter interface which provides methods for setting
and retrieving a standard and error writer. You pass this value to the
pager.Start method and if either of the writers is connected to a terminal
then it will start a pager command and insert it between the writer and the
terminal. The Start method returns a pager and you must call the Done method
on the returned pager.

*/
package pager
