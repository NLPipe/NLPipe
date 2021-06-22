args=("$@")

# parsing args from command line
full_file_name="${args[0]}"
sentences_max_num=${args[1]}

if [[ -z "${sentences_max_num// }" ]]; then
  sentences_max_num=1
fi

i=1

echo "Removing old files in current ../sentences/ directory..."
rm -f ../sentences/*.txt
echo "Old files removal DONE!"

while read line; do
  line_escaped=$(echo $line | sed 's/"/\\"/g')

  echo $line_escaped >> "../sentences/${i}.txt"

  echo "Sentence $i / $sentences_max_num DONE!"

  ((i=i+1))

  if (( $i > $sentences_max_num )); then
    break
  fi
done < $full_file_name
