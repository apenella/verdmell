{{ define "scripts" }}

	generateServicesDetails = 

	$(document).ready(function(){
		$(".clusternode").click(function () {
			var url = $(this).attr('url')+"/api/services"
			$.getJSON(url, function(services) {
				var servicesinfo = '<div class="servicesdetail">';
				$.each(services.services, function(key, service) {
					servicesinfo += '<div id="'+service.name+'" class="servicedetail">';
						servicesinfo += '<div class="servicename">';
						servicesinfo += service.name;
						servicesinfo += '</div>';					
						servicesinfo += '<div class="servicechecks">';
							for ( c in service.checks) {
								servicesinfo += '<div class="servicecheck">'+service.checks[c]+'</div>'
							}
						servicesinfo += '</div>';
					servicesinfo += '</div>';			
				});
				servicesinfo += '</div>';
				
				$(".clusternodedetails").html(servicesinfo);
			});

		});

	});
{{end}}