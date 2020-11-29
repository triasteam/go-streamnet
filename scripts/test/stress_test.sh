for ((i = 0; i <= 255; i++)) 
do
	node1=192.168.0.$i
	j=$((i+1))
	node2=192.168.0.$j
	#curl -X POST -d '{"Attester": "'"$node1"', "Attestee": "'"$node2"', "Score": "1"}' http://127.0.0.1:14700/save
	curl -X POST -d "{\"Attester\": \"$node1\", \"Attestee\": \"$node2\", \"Score\": \"1\"}" http://127.0.0.1:14700/save
done
