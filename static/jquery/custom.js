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
		var lastfmUser = $('#lastfmUser').val();
		window.location.href = "/sched/" + lastfmUser + "/" + selectedFestival.name;
	}
	else {
		alert("Please select a festival");
	}
});

