import curses

myscreen = curses.initscr()

myscreen.border(0)
myscreen.addstr(12, 25, "Python curses is action!")
myscreen.refresh()
myscreen.getch()
curses.endwin()