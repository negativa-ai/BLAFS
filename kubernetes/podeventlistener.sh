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
	 echo  "## processing image $image in pod/$line with owner $owner"
         ( echo $image | grep baffs ) && echo ignored && return 
         cat writepipe.yaml | sed s/CONTAINER/$image/ | sed s/NAME/writepipe-$RANDOM/| kubectl create -f -
      done

      sleep 20
      echo patching $owner
      kubectl get $owner -o yaml |grep -v -E 'image:.*-baffs$' | sed 's/image:.*$/&-baffs/' | kubectl apply -f -
      echo deleting pod $line
      while ! kubectl delete pod $line 
      do
	      if  kubectl get pods
	      then 
		      break
	      fi
	      echo Retry in 5 seconds
	      sleep 5
      done
    fi
  fi

}

while true
do
    kubectl get event -w --field-selector reason=Started,involvedObject.kind=Pod \
       	-o go-template='{{printf "%s\n" .involvedObject.name}}' | \
    while read line
    do
      debloat_cycle $line # &
      kubectl get pods -o go-template --template='{{"\nImages\n"}}
    {{- range .items -}}
      {{- printf "%s\n" .metadata.name -}}
      {{- range .spec.initContainers -}}
       {{"*  "}} {{ .image }}
    {{- end -}}
    {{- range .spec.containers -}}
       {{"-  "}} {{ .image }}
    {{- end -}}{{"\n"}}
{{- end -}}' | grep --color=always -E '^.*-baffs|$'

    done
    sleep 10
done

