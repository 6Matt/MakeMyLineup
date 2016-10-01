
var festivals = new Bloodhound({
  datumTokenizer: Bloodhound.tokenizers.whitespace,
  queryTokenizer: Bloodhound.tokenizers.whitespace,
  prefetch: '/data/festivalNames.json'
});

$('#festivalInput').typeahead(null, {
  name: 'festivals',
  source: festivals
});

$('#festivalInput').keyup(function(event){
	if(event.keyCode == 13) {
		alert($('#festivalInput').val()); //TODO - change this to call the schedule maker
	}
});