import java.util.List;
import java.util.Collections;
import java.util.ArrayList;
import java.nio.charset.StandardCharsets;
import java.io.PrintStream;
import java.io.IOException;

public class Combination{
    static void swap(List<Integer> a, int i, int j) {
        Collections.swap(a,i,j);
    }

    static String toString(List<Integer> a) {
        String s = ""; 
        for(int codePoint : a){
            s += Character.toString((char)codePoint);
        }
        return s;
    }

    static void all(List<Integer> s, int n) {
        if (n == 1) {
            /* try{
               PrintStream out = new PrintStream(System.out, true, "UTF-8");
               out.println(toString(s));
               }catch (IOException e){
               } 
               */
        } else {
            for (int i = 0; i < n; i++) {
                swap(s, i, n - 1);
                all(s, n - 1);
                swap(s, i, n - 1);
            }
        }
    }
    public static void main(String []args) {
        String a = "♠行到水穷处，坐看云起时";
        List<Integer> b = new ArrayList<Integer>();
        for(int i = 0;i < a.length();i++)  {
            b.add(a.codePointAt(i));
        }
        all(b, b.size());
    }
}
