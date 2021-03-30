CREATE DATABASE netmap;
CREATE EXTENSION IF NOT EXISTS timescaledb;

\c netmap

CREATE TABLE graphs (
  time        TIMESTAMPTZ       NOT NULL,
  reporter    TEXT              NOT NULL,
  vertices    JSON              NULL,
  edges       JSON              NULL,
  packets     INTEGER           NULL
);

SELECT create_hypertable('graphs', 'time');
SELECT add_retention_policy('graphs', INTERVAL '30 minutes');


create or replace view edges_1min as 
select 
	agg."time",
	json_agg(
		json_build_object(
			'source', agg.sourceIp,
			'destination', agg.destinationIp,
			'properties', json_build_object(
				'trafficType', agg.trafficType,
				'weight', agg.weight,
				'sourcePort', agg.sourcePort,
				'destinationPort', agg.destinationPort,
				'packetCount', agg.packetCount
			)
		) 
	) as obj,
	sum(agg.packetCount)
from (
	select 
		time_bucket(interval '1 minute', g."time") as "time",
		obj->>'source' as sourceIp, 
		cast(obj->'properties'->>'sourcePort' as integer) as sourcePort, 
		obj->>'destination' as destinationIp, 
		cast(obj->'properties'->>'destinationPort' as integer) as destinationPort, 
		obj->'properties'->>'trafficType' as trafficType, 
		avg(cast(obj->'properties'->>'weight' as double precision)) as weight, 
		sum(cast(obj->'properties'->>'packetCount' as integer)) as packetCount
	from graphs g, json_array_elements(g.edges) obj
	group by time_bucket(interval '1 minute', g."time"), sourceIp, sourcePort, destinationIp, destinationPort, trafficType
) as agg
group by "time"
order by "time";

create or replace view vertices as 
select distinct on (obj->>'id') obj->>'id' as hostId, obj
from graphs g, json_array_elements(g.vertices) obj;


select jsonb_agg(unique_vertices.obj) as vertices
from (
	select distinct on (obj->>'id') obj->>'id' as hostId, obj
	from graphs g, json_array_elements(g.vertices) obj
) unique_vertices

select jsonb_agg(agg_edges.obj) as vertices
from (
	select json_build_object(
		'source', agg.sourceIp,
		'destination', agg.destinationIp,
		'properties', json_build_object(
			'trafficType', agg.trafficType,
			'weight', agg.weight,
			'sourcePort', agg.sourcePort,
			'destinationPort', agg.destinationPort,
			'packetCount', agg.packetCount
		) 
	) as obj
	from (
		select 
		obj->>'source' as sourceIp, 
		obj->'properties'->>'sourcePort' as sourcePort, 
		obj->>'destination' as destinationIp, 
		obj->'properties'->>'destinationPort' as destinationPort, 
		obj->'properties'->>'trafficType' as trafficType, 
		avg(cast(obj->'properties'->>'weight' as double precision)) as weight, 
		sum(cast(obj->'properties'->>'packetCount' as integer)) as packetCount
		from graphs g, json_array_elements(g.edges) obj
		group by sourceIp, sourcePort, destinationIp, destinationPort, trafficType
	) as agg
	
) agg_edges

