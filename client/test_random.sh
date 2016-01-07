long=1
short=1
echo $end

while ([ $long -gt 0 ] || [ $short -gt 0 ] )
do
	num=num=$(($RANDOM+1000000000))
	rand=$(($num%10+1))
	if ([ $rand -ge 5 ] && [ $long -gt 0 ])
	then
		echo "run long job"
		long=$(($long-1)) 
		./client subjob sub_1.json
	elif ([ $rand -lt 5 ] && [ $short -gt 0 ])
	then
		echo "run short job"
		short=$(($short-1))
		./client subjob sub_2.json
	fi
done

	
