#include <stdio.h>
#include <stdint.h>
#include <unistd.h>

 void itoa(int n)
 {
     int i, sign;
     char s[64];
 
     if ((sign = n) < 0)  /* record sign */
         n = -n;          /* make n positive */
      i = 0;
     s[sizeof(s) - 1] = '\n';
     i++;
     do {       /* generate digits in reverse order */
         s[sizeof(s) - i - 1] = n % 10 + '0';   /* get next digit */
         i++;
     } while ((n /= 10) > 0);     /* delete it */

     if (sign < 0) {
         s[sizeof(s) - i - 1] = '-';
         i++;
     }

     s[sizeof(s) - i - 1] = '\0';
     i++;
     write(1, &s[sizeof(s) - i], i);
 }

void main() {
    itoa(100);
    itoa(-100);
    itoa(-0);
}
