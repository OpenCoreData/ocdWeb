Slim.tag(
	'repeater-child', 
	'<div class="div-table-row"> \
	     <div class="div-table-colsmall" bind:prop="item.value" bind>{{item.value}}</div> \
	    <div class="div-table-colsmall" bind:prop="item.unitText" bind>{{item.unitText}}</div>\
	    <div class="div-table-col" bind:prop="item.description" bind>{{item.description}}</div>\
	 </div>', class extends Slim { });
Slim.tag(
	'geocomponents-parameters',
	`<div class="div-table">
	           <repeater-child s:repeat="items as item"></repeater-child> 
    </div>`,
	class ParamTag extends Slim {
		// your code here
		onBeforeCreated() {
			var element = document.getElementById('schemaorg');
			var jsonld = element.innerHTML;
			var obj = JSON.parse(jsonld);
			this.items = obj.variableMeasured
		}
		myMethod() {
			return "test"
		}
	}
)
