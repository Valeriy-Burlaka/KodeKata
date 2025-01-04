#include <stdio.h>
#include <termios.h>
#include <unistd.h>

// Demonstrate what happens when we turn off an automatic conversion of "\n" to "\r\n" on Linux systems
int main() {
    struct termios term;
    tcgetattr(STDOUT_FILENO, &term);

    // Turn off automatic \r insertion
    term.c_oflag &= ~ONLCR;
    tcsetattr(STDOUT_FILENO, TCSANOW, &term);

    printf("First line\nSecond line\n");  // Now this will print diagonally!

    // Reset to normal
    term.c_oflag |= ONLCR;
    tcsetattr(STDOUT_FILENO, TCSANOW, &term);
}
