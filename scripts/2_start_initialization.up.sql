WITH json (doc) AS (
   values
       ('[
         {
           "country_name": "Australia",
           "picture": "."
         },
         {
           "country_name": "Armenia",
           "picture": "."
         },
         {
           "country_name": "Belarus",
           "picture": "."
         },
         {
           "country_name": "Brazil",
           "picture": "."
         },
         {
           "country_name": "Canada",
           "picture": "."
         },
         {
           "country_name": "France",
           "picture": "."
         },
         {
           "country_name": "Germany",
           "picture": "."
         },
         {
           "country_name": "India",
           "picture": "."
         },
         {
           "country_name": "Japan",
           "picture": "."
         },
         {
           "country_name": "Republic_of_China",
           "picture": "."
         },
         {
           "country_name": "Russia",
           "picture": "."
         },
         {
           "country_name": "Spain",
           "picture": "."
         },
         {
           "country_name": "Switzerland",
           "picture": "."
         },
         {
           "country_name": "Turkey",
           "picture": "."
         },
         {
           "country_name": "UK",
           "picture": "."
         },
         {
           "country_name": "Ukraine",
           "picture": "."
         },
         {
           "country_name": "USA",
           "picture": "."
         },
         {
           "country_name": "Vietnam",
           "picture": "."
         }
       ]'::json)
)
INSERT INTO countries (country_name, picture)
SELECT p.country_name, p.picture
from json l cross join lateral json_populate_recordset(null::countries, doc) as p;


WITH json (doc) AS (
    values
        ('[
          {
            "language": "Chinese",
            "picture": "."
          },
          {
            "language": "English",
            "picture": "."
          },
          {
            "language": "Spanish",
            "picture": "."
          },
          {
            "language": "Arabian",
            "picture": "."
          },
          {
            "language": "Russian",
            "picture": "."
          },
          {
            "language": "Portuguese",
            "picture": "."
          },
          {
            "language": "French",
            "picture": "."
          },
          {
            "language": "German",
            "picture": "."
          },
          {
            "language": "Hindi",
            "picture": "."
          },
          {
            "language": "Bengali",
            "picture": "."
          },
          {
            "language": "Japanese",
            "picture": "."
          },
          {
            "language": "Italian",
            "picture": "."
          },
          {
            "language": "Belarusian",
            "picture": "."
          },
          {
            "language": "Ukrainian",
            "picture": "."
          },
          {
            "language": "Vietnamese",
            "picture": "."
          },
          {
            "language": "Armenian",
            "picture": "."
          },
          {
            "language": "Turkish",
            "picture": "."
          }
        ]'::json)
)
INSERT INTO languages (language, picture)
SELECT p.language, p.picture
FROM json l CROSS JOIN lateral json_populate_recordset(null::languages, doc) AS p;

