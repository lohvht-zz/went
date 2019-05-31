function Point(x, y) {
  this.x = x;
  this.y = y;
}
var p1 = new Point(11, 22);
var p2 = new Point(33, 44);
console.log(p1.size); // undefined
Point.prototype.size = 100;
console.log(p1.size) // 100
console.log(p2.size); // 100
p1.__proto__ === Point.prototype; // true



var p1 = new Point(11, 22);
function getX(p) { return p.x; }
var sumX = 0;
sumX += getX(p1);