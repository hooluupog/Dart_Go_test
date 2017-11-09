/* 
 * The difference and usages of 'abstract class','mixin' and 'interface'.
 */

/// Abstract class. Can be implemented ,some-implemented or
/// unimplemented but can not be  instantiated.
abstract class Human {
  void profile();
  String major() => 'Unknown';
}

abstract class Philosopher {
  String thinking() => 'Thinking the nature of the world.';
  String major() => 'Philosopher';
}

abstract class Musician {
  String create() => 'Creating good music.';
  String major() => 'Musician';
}

abstract class Athlete {
  String compete() => 'For the champion.';
  String major() => 'Athlete';
}

/// Mixins. 'included' relationship.
class Person extends Human with Philosopher, Musician, Athlete {
  String _name;
  Person(this._name);
  @override
  profile() => print('''
  name: $_name
  major: ${major()}
  hobbies:
    ${Hobbyies.philosopher(this)}
    ${Hobbyies.musician(this)}
    ${Hobbyies.athlete(this)}''');
}

/// Interface. Subclass does not inherit but implement(must implement)
/// other classes' methods.It's a 'has a' relationship.
class Student extends Human implements Philosopher, Musician, Athlete {
  String _name;
  Student(this._name);
  String thinking() => '''Asking 'what','why', and 'how'.''';
  String create() => 'Learning the music.';
  String compete() => 'Playing game.';
  @override
  profile() => print('''
  name: $_name
  major: ${major()}
  hobbies:
    ${Hobbyies.philosopher(this)}
    ${Hobbyies.musician(this)}
    ${Hobbyies.athlete(this)}''');
}

/// Interface caller.
class Hobbyies {
  static String philosopher(Philosopher p) => p.thinking();
  static String musician(Musician m) => m.create();
  static String athlete(Athlete a) => a.compete();
}

main() {
  var person = new Person('Bob');
  person.profile();
  var stu = new Student('Tom');
  stu.profile();
}
