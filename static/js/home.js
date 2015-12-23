var connection = new WebSocket('ws://'+location.host+'/ws');

var map;
function initMap() {
	navigator.geolocation.getCurrentPosition(function(position) {
		a = position.coords.latitude; b = position.coords.longitude;

		map = new google.maps.Map(document.getElementById('map'), {
			center: {lat: a, lng: b},
			zoom: 14
		});
	});

}


connection.onmessage = function(event){

	//console.log(event.data);


	var drop = JSON.parse(event.data);
	$("#droplist").prepend($('<li class="dropitem">').text(drop["name"] + ": " + drop["message"]));
	new google.maps.Marker({
		position: {lat: drop["y"], lng: drop["x"]},
		map: map,
		title:drop["name"]
	});

}



jQuery(document).ready(function($) {
	
	$("#dataclick").on('click', function(event) {

		var max = 0;
		var happy = 0;
		var sad = 0;


		$( "li" ).each(function( index ) {
			//console.log( index + ": " + $( this ).text() );
			max = index;

			value = $(this).text();


			if(value.search("happy") > 0){
				happy++;
			}
			if(value.search("great") > 0){
				happy++;
			}
			if(value.search("love") > 0){
				happy++;
			}
			if(value.search("amazing") > 0){
				happy++;
			}
			if(value.search("fun") > 0){
				happy++;
			}
			if(value.search("nice") > 0){
				happy++;
			}

			if(value.search("sad") > 0){
				sad++;
			}
			if(value.search("hate") > 0){
				sad++;
			}
			if(value.search("not like") > 0){
				sad++;
			}
			if(value.search("bad") > 0){
				sad++;
			}
			if(value.search("mean") > 0){
				sad++;
			}
			//unhappy
		});

		console.log(sad);

		$("#happy").text(happy/max);
		$("#unhappy").text(sad/max);

	});


	navigator.geolocation.getCurrentPosition(function(position) {
		a = position.coords.latitude; b = position.coords.longitude;

		//onsole.log("http://"+location.host+"/get/?r=10&x="+b+"&y="+a);

		$.ajax({
			url: "http://"+location.host+"/get/?r=10&x="+b+"&y="+a
		})
		.done(function(data) {
			 //console.log(data)
			 for (var i = 0; i < data.length;i++){
				//console.log(data[i])
				$("#droplist").append($('<li class="dropitem"></li>').text(data[i]["name"] + ": " + data[i]["message"]));


				new google.maps.Marker({
					position: {lat: data[i]["y"], lng: data[i]["x"]},
					map: map,
					title: data[i]["name"]
				});



			}
		});

	});
	
	$("#dropit").on("click", function(){

		var date = new Date()
		datestring = date.getFullYear() +"-"+ (date.getMonth()+1) +"-"+ date.getDate();

		var hour= 0;
		if (date.getHours() < 10){
			hour = "0" + date.getHours();
		} else {
			hour = date.getHours();
		}

		var minute = 0;
		if (date.getMinutes() < 10){
			minute = "0" + date.getMinutes();
		} else {
			minute = date.getMinutes();
		}

		timestring = hour + ":" + minute;

		var a; var b;

		navigator.geolocation.getCurrentPosition(function(position) {
			a = position.coords.latitude; b = position.coords.longitude;

			var data = {
				name: $("#name").val(),
				message: $("#message").val(),
				x: b,
				y: a,
				date: datestring,
				time: timestring
			}

			//console.log(JSON.stringify(data));
			connection.send(JSON.stringify(data));

			$("#name").val("");
			$("#name").val("");


		});

		
		
	});
	
	


	$('.modal-trigger').leanModal();



});