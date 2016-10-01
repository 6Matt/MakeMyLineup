
var festivals = new Bloodhound({
  datumTokenizer: Bloodhound.tokenizers.obj.whitespace('desc'),
  queryTokenizer: Bloodhound.tokenizers.whitespace,
  prefetch: '/data/festivals.json'
});

$('#festivalInput').typeahead(null, {
  name: 'festivals',
  display: 'desc',
  source: festivals
});

$('#festivalInput').keyup(function(event){
	if(event.keyCode == 13) {
		alert($('#festivalInput').val()); //TODO - change this to call the schedule maker
	}
});