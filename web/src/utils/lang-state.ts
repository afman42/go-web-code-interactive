import { writable } from "svelte/store";
export let langState = writable({
  value: "node",
  sampleDataLang: {
    repl: {
      node: "console.log('hello world')\n\n\n\n\n\n",
      php: "<?php\necho 'hello world';\n\n\n\n\n",
      go: 'package main\n\nimport "fmt"\n\nfunc main(){\n\tfmt.Println("hello world")\n}\n\n\n\n\n',
    },
    stq: {
      node: "//Don't remove function intIntoString\n\nfunction intIntoString(n){\n\treturn\n}\n\n",
      php: "<?php\n//Don't remove function intIntoString\n\nfunction intIntoString($n){\n\treturn\n}\n\n",
      go: 'package main\n\nimport "fmt" //Don"t remove fmt\n\n//Don"t remove function intIntoString\nfunc intIntoString(n int) string {\n\treturn\n}\n\n',
    },
  },
  type: "repl",
});
