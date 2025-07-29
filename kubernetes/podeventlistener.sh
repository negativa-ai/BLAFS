#!/bin/bash
#
#
# gets the original owner of a pod
#

#set -x

cleanup() {
    # kill all processes whose parent is this process
    pkill -P $$
}

for sig in INT QUIT HUP TERM; do
  trap "
    cleanup
    trap - $sig EXIT
    kill -s $sig "'"$$"' "$sig"
done
trap cleanup EXIT



getowner(){
  local object=$1
  local owner=$(kubectl get $object -o jsonpath='{ .metadata.ownerReferences[].kind}/{.metadata.ownerReferences[].name}' || return 1)
 

  if [ -z $owner ]
  then
    return 1
  fi

  if [ $owner == '/' ] 
  then 
    echo $object
    return 0
  else
    getowner $owner
    return $?
  fi
}

debloat_cycle(){
  local line=$1
  if [ ! -z $line ]
  then
    owner=$(getowner pod/$line)
    if [[ ! -z $owner && ! $owner == / ]]
    then
      for image in $(kubectl get pod $line  -o jsonpath="{ .spec['initContainers', 'containers'][*].image}
            ")
      do 
	 echo -n "## "
         ( echo $image | grep baffs ) && echo ignored && return 
         cat writepipe.yaml | sed s/CONTAINER/$image/ | sed s/NAME/writepipe-$RANDOM/| kubectl create -f -
      done

      sleep 20
      if  kubectl delete pod $line
      then
          sleep 20
	  echo patching $owner
          kubectl get $owner -o yaml |grep -v -E 'image:.*-baffs$' | sed 's/image:.*$/&-baffs/' | kubectl apply -f -
      fi
    fi
  fi

}

kubectl get event -w --field-selector reason=Started,involvedObject.kind=Pod \
       	-o go-template='{{printf "%s\n" .involvedObject.name}}' |
(
while true
do
  read line 
  debloat_cycle $line # &
done
)
