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

var selectedFestival;

$('#festivalInput').bind('typeahead:select', function(ev, suggestion) {
	selectedFestival = suggestion;
});
$('#festivalInput').bind('typeahead:autocomplete', function(ev, suggestion) {
	selectedFestival = suggestion;
});

$('#makeSchedule').click(function(){
	if(typeof selectedFestival != "undefined" && $('#festivalInput').val() == selectedFestival.desc) {
		alert(selectedFestival.name); //TODO - change this to call the schedule maker
	}
	else {
		alert("Select a festival");
	}
});

