package utils

// Todo need test case
// Need Refactor
const TxtGo = `
func main(){
  data := true
  f := intIntoString(1)
  _, err := strconv.Atoi(f)
  if err != nil{
    fmt.Println(err)
  } 

  data = false
  fmt.Println(data)
}
`

const TxtJS = `
function main(f){
  let data = true
  if(parseInt(f,10).toString()===f) {
    data = false
  }
  return data;
}
console.log(main(intIntoString(1)))
`

const TxtPHP = `
function main($f){
  $data = true;
  if(ctype_digit($f)) {
    $data = false;
  }
  return $data ? "true" : "false";
}
echo main(intIntoString(1));
`

// Todo
// try {
//   echo main(intIntoString(1));
// }
// //catch exception
// catch(Exception $e) {
//   echo 'Message: ' .$e->getMessage();
// }
