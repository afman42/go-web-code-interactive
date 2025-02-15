export let langState = $state({
  value: "node",
  sampleDataLang: {
    repl: {
      node: "console.log('hello world')\n\n\n\n\n\n",
      php: "<?php\necho 'hello world';\n\n\n\n\n",
      go: 'package main\n\nimport "fmt"\n\nfunc main(){\n\tfmt.Println("hello world")\n}\n\n\n\n\n',
    },
    stq: {
      node: "//Don't remove function intIntoString\nfunction intIntoString(n){\n\treturn\n}\n\n",
      php: "<?php\n\n//Don't remove function intIntoString\nfunction intIntoString($n){\n\treturn\n}\n\n",
      go: 'package main\n\nimport (\n\t"fmt" //Don"t remove\n\t"strconv"  //Don"t remove\n)\n\n//Don"t remove function intIntoString\nfunc intIntoString(n int) string {\n\treturn\n}\n\n',
    },
  },
  type: "repl",
}) as {
  value: string;
  sampleDataLang: Record<string, Record<string, string>>;
  type: string;
};
